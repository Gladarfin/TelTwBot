package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Stats struct {
	ID           int
	Name         string
	DisplayName  string
	MinValue     int
	MaxValue     int
	DefaultValue int
	CreatedAt    sql.NullTime
}

type UserStats struct {
	UserID    int
	StatType  string
	Value     int
	UpdatedAt sql.NullTime
}

func (d *Database) GetUserStats(ctx context.Context, userID int) ([]UserStats, error) {
	const query = `
			SELECT s.name, us.value, us.updated_at
			FROM user_stats us
			JOIN stat_types s ON us.stat_type_id = s.id
			WHERE us.user_id = $1
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user stats: %s", err)
	}
	defer rows.Close()

	var stats []UserStats
	for rows.Next() {
		var stat UserStats
		stat.UserID = userID
		if err := rows.Scan(&stat.StatType, &stat.Value, &stat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Failed to scan stat row: %s", err)
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

func (d *Database) UpdateUserStats(ctx context.Context, userID int, statName string, value int) error {

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var stat Stats
	const queryStatConstraints = `
			SELECT id, min_value, max_value
			from stat_types
			WHERE name = $1
	`
	err = tx.QueryRowContext(ctx, queryStatConstraints, statName).Scan(&stat.ID, &stat.MinValue, &stat.MaxValue)
	if err != nil {
		return fmt.Errorf("Failed to get stat type constraints: %w", err)
	}

	var curValue int
	const queryCurrentStatValue = `
			SELECT value
			FROM user_stats
			WHERE user_id = $1 AND stat_type_id = $2
	`
	err = tx.QueryRowContext(ctx, queryCurrentStatValue, userID, stat.ID).Scan(&curValue)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Failed to get current stat value: %w", err)
	}

	if err == sql.ErrNoRows {
		curValue = stat.DefaultValue
	}

	newValue := curValue + value
	if newValue < stat.MinValue {
		newValue = stat.MinValue
	}
	if newValue > stat.MaxValue {
		newValue = stat.MaxValue
	}

	const queryUpdateValue = `
			INSERT INTO user_stats (user_id, stat_type_id, value)
			VALUES($1, $2, $3)
			ON CONFLICT (user_id, stat_type_id) DO UPDATE
			SET value = EXCLUDED.value, updated_at = NOW()
	`

	_, err = tx.ExecContext(ctx, queryUpdateValue, userID, stat.ID, newValue)
	if err != nil {
		return fmt.Errorf("Failed to update stat: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}
