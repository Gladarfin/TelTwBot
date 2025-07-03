package database

import (
	"context"
	"database/sql"
	"fmt"
)

type UserResult struct {
	UserID    int
	Wins      int
	Draws     int
	Loses     int
	UpdatedAt sql.NullTime
}

func (d *Database) GetUserResults(ctx context.Context, userID int) (*UserResult, error) {
	const query = `
			SELECT total_wins, total_draws, total_lose, updated_at
			FROM user_results
			WHERE user_id = $1
	`

	var result UserResult
	result.UserID = userID
	err := d.db.QueryRowContext(ctx, query, userID).Scan(
		&result.Wins,
		&result.Draws,
		&result.Loses,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user results: %w", err)
	}

	return &result, nil
}

func (d *Database) IncrementUserResult(ctx context.Context, userID int, resultType string) error {
	var field string

	switch resultType {
	case "win":
		field = "total_wins"
	case "draw":
		field = "total_draws"
	case "lose":
		field = "total_lose"
	default:
		return fmt.Errorf("invalid result type: %s", resultType)
	}

	query := fmt.Sprintf(`
		INSERT INTO user_results (user_id, %s)
		VALUES ($1, 1)
		ON CONFLICT (user_id) DO UPDATE
		SET %s = user_results.%s + 1, updated_at = NOW()
	`,
		field, field, field)

	_, err := d.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment user results: %s", err)
	}

	return nil
}
