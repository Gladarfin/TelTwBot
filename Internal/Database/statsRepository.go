package database

import (
	"context"
	"database/sql"
	"errors"
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

func (d *Database) getOrCreateUser(ctx context.Context, username string) (int, error) {
	var userID int
	const userIdQuery = `
			SELECT id FROM users WHERE username = $1
	`
	err := d.db.QueryRowContext(ctx, userIdQuery, username).Scan(&userID)
	if err == nil {
		return userID, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("error checking user: %w", err)
	}

	const addUserQuery = `
			INSERT INTO users (username, created_at, updated_at)
			VALUES($1, NOW(), NOW())
			RETURNING id
	`
	err = d.db.QueryRowContext(ctx, addUserQuery, username).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	return userID, nil
}

func (d *Database) GetOrCreateUserStats(ctx context.Context, username string) ([]UserStats, error) {

	userID, err := d.getOrCreateUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user setup failed: %w", err)
	}

	if err := d.createDefaultStatsForUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to ensure stats exist: %w", err)
	}

	return d.getExistingUserStats(ctx, userID)
}

func (d *Database) getExistingUserStats(ctx context.Context, userID int) ([]UserStats, error) {
	const query = `
			SELECT s.name, us.value, us.updated_at
			FROM user_stats us
			JOIN stat_types s ON us.stat_type_id = s.id
			WHERE us.user_id = $1
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	defer rows.Close()

	var stats []UserStats
	for rows.Next() {
		var stat UserStats
		stat.UserID = userID
		if err := rows.Scan(&stat.StatType, &stat.Value, &stat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stat row: %w", err)
		}

		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return stats, nil
}

func (d *Database) createDefaultStatsForUser(ctx context.Context, userID int) error {

	statTypes, err := d.getAllStatTypes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get stat types: %w", err)
	}

	const insertMissingStats = `
			INSERT INTO user_stats (user_id, stat_type_id, value, updated_at)
			SELECT $1, id, default_value, NOW()
			FROM stat_types
			WHERE name = $2
			ON CONFLICT (user_id, stat_type_id) DO NOTHING
	`
	for _, statType := range statTypes {
		_, err := d.db.ExecContext(ctx, insertMissingStats, userID, statType)
		if err != nil {
			return fmt.Errorf("failed to create stat %s: for user: %w", statType, err)
		}
	}

	return nil
}

func (d *Database) UpdateUserStats(ctx context.Context, userID int, statName string, value int) error {

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
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
		return fmt.Errorf("failed to get stat type constraints: %w", err)
	}

	var curValue int
	const queryCurrentStatValue = `
			SELECT value
			FROM user_stats
			WHERE user_id = $1 AND stat_type_id = $2
	`
	err = tx.QueryRowContext(ctx, queryCurrentStatValue, userID, stat.ID).Scan(&curValue)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get current stat value: %w", err)
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
		return fmt.Errorf("failed to update stat: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (d *Database) getAllStatTypes(ctx context.Context) ([]string, error) {
	const query = `SELECT name FROM stat_types`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get stat type: %w", err)
	}
	defer rows.Close()

	var statTypes []string
	for rows.Next() {
		var statType string
		if err := rows.Scan(&statType); err != nil {
			return nil, fmt.Errorf("failed to scan stat type: %w", err)
		}
		statTypes = append(statTypes, statType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return statTypes, nil
}

func (d *Database) UpdateUserStat(ctx context.Context, username string, stat string, val int) (string, int, error) {

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to begin transaction %w", err)
	}

	defer tx.Rollback()

	var userID int
	err = tx.QueryRowContext(ctx, "SELECT id from users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("user %s not found", username)
		}
		return "", 0, fmt.Errorf("failed to get user: %w", err)
	}

	var statTypeID, minValue, maxValue int
	var statName string
	err = tx.QueryRowContext(ctx, `
		SELECT id, name, min_value, max_value
		FROM stat_types
		WHERE name = $1`, stat).Scan(&statTypeID, &statName, &minValue, &maxValue)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("stat %s not found", stat)
		}
		return "", 0, fmt.Errorf("failed to get stat type: %w", err)
	}

	var curStatValue, curFreePoints int
	err = tx.QueryRowContext(ctx, `
		SELECT value FROM user_stats
		WHERE user_id = $1 AND stat_type_id = $2
	`, userID, statTypeID).Scan(&curStatValue)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get current stat value: %w", err)
	}

	var freePointsStatID int
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM stat_types
		WHERE name = 'free-points'
	`).Scan(&freePointsStatID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get free-points stat id: %w", err)
	}

	err = tx.QueryRowContext(ctx, `
		SELECT value FROM user_stats
		WHERE user_id = $1 AND stat_type_id = %2
	`, userID, freePointsStatID).Scan(&curFreePoints)

	if err != nil {
		if err == sql.ErrNoRows {
			curFreePoints = 0
		} else {
			return "", 0, fmt.Errorf("failed to get free points: %w", err)
		}
	}

	if curFreePoints < val {
		return "", 0, fmt.Errorf("not enough free points: need %d, have %d", val, curFreePoints)
	}

	newStatValue := curStatValue + val
	if newStatValue > maxValue {
		newStatValue = maxValue
	}

	pointsUsed := newStatValue - curStatValue
	newFreePoints := curFreePoints - pointsUsed

	const updStatQuery = `
		INSERT INTO user_stats (user_id, stat_type_id, value)
		VALUES ($1,$2,$3)
		ON CONFLICT (user_id, stat_type_id)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`
	_, err = tx.ExecContext(ctx, updStatQuery, userID, statTypeID, newStatValue)

	if err != nil {
		return "", 0, fmt.Errorf("failed to update stat: %w", err)
	}

	_, err = tx.ExecContext(ctx, updStatQuery, userID, freePointsStatID, newFreePoints)

	if err != nil {
		return "", 0, fmt.Errorf("failed to update free points: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return statName, newStatValue, nil
}
