package database

import (
	"context"
	"database/sql"
	"fmt"
	"math"
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

func (d *Database) UpdateResultsAfterDuel(ctx context.Context, initiator string, challenger string, result int) error {
	tx, err := d.db.(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	tempRepo := &Database{db: tx}

	challengerId, err := tempRepo.getUserIdByUsername(ctx, challenger)
	if err != nil {
		return fmt.Errorf("failed to get challenger ID: %w", err)
	}

	initiatorId, err := tempRepo.getUserIdByUsername(ctx, initiator)
	if err != nil {
		return fmt.Errorf("failed to get initiator ID: %w", err)
	}

	switch result {
	case 0:
		if err := tempRepo.IncrementUserResult(ctx, initiatorId, "draw"); err != nil {
			return fmt.Errorf("failed to update initiator result values: %w", err)
		}
		if err := tempRepo.IncrementUserResult(ctx, challengerId, "draw"); err != nil {
			return fmt.Errorf("failed to update challenger result values: %w", err)
		}
		if err := tempRepo.updateUserPoints(ctx, initiatorId, "draw"); err != nil {
			return fmt.Errorf("failed to update initiator stats: %w", err)
		}
		if err := tempRepo.updateUserPoints(ctx, challengerId, "draw"); err != nil {
			return fmt.Errorf("failed to update challenger stats: %w", err)
		}
	case 1:
		if err := tempRepo.IncrementUserResult(ctx, initiatorId, "win"); err != nil {
			return fmt.Errorf("failed to update initiator result values: %w", err)
		}
		if err := tempRepo.IncrementUserResult(ctx, challengerId, "lose"); err != nil {
			return fmt.Errorf("failed to update challenger result values: %w", err)
		}
		if err := tempRepo.updateUserPoints(ctx, initiatorId, "win"); err != nil {
			return fmt.Errorf("failed to update initiator stats: %w", err)
		}
	case 2:
		if err := tempRepo.IncrementUserResult(ctx, initiatorId, "lose"); err != nil {
			return fmt.Errorf("failed to update initiator result values: %w", err)
		}
		if err := tempRepo.IncrementUserResult(ctx, challengerId, "win"); err != nil {
			return fmt.Errorf("failed to update challenger result values: %w", err)
		}
		if err := tempRepo.updateUserPoints(ctx, challengerId, "win"); err != nil {
			return fmt.Errorf("failed to update challenger stats: %w", err)
		}
	default:
		return fmt.Errorf("invalid game result: %d", result)
	}

	return tx.Commit()
}

func (d *Database) getUserIdByUsername(ctx context.Context, username string) (int, error) {
	const getUserIdQuery = `
			SELECT id FROM users WHERE username = $1
	`
	var id int
	err := d.db.QueryRowContext(ctx, getUserIdQuery, username).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}
	return id, nil
}

func (d *Database) updateUserPoints(ctx context.Context, userID int, result string) error {
	var winsCount, drawsCount int
	err := d.db.
		QueryRowContext(ctx, `SELECT total_wins, total_draws FROM user_results WHERE user_id = $1`, userID).
		Scan(&winsCount, &drawsCount)
	if err != nil {
		return fmt.Errorf("couldn't get user points: %w", err)
	}

	const getStatIdQuery = `SELECT id FROM stat_types WHERE name = $1`
	var statTypeIDTotalFreePoints, statTypeIDFreePoints int

	err = d.db.QueryRowContext(ctx, getStatIdQuery, "free-points").Scan(&statTypeIDFreePoints)
	if err != nil {
		return fmt.Errorf("couldn't get user free-points stat: %w", err)
	}

	err = d.db.QueryRowContext(ctx, getStatIdQuery, "total-free-points").Scan(&statTypeIDTotalFreePoints)
	if err != nil {
		return fmt.Errorf("couldn't get user totat free points stat: %w", err)
	}

	const getStatTypeIdQuery = `SELECT value FROM user_stats WHERE user_id = $1 AND stat_type_id = $2`
	var freePoints, totalFreePoints int
	err = d.db.QueryRowContext(ctx, getStatTypeIdQuery, userID, statTypeIDFreePoints).Scan(&freePoints)
	if err != nil {
		return fmt.Errorf("couldn't get user free points value: %w", err)
	}
	err = d.db.QueryRowContext(ctx, getStatTypeIdQuery, userID, statTypeIDTotalFreePoints).Scan(&totalFreePoints)
	if err != nil {
		return fmt.Errorf("couldn't get user total free points value: %w", err)
	}

	winThreshold := 10 + 10*int(math.Floor(float64(totalFreePoints)/5))
	drawThreshold := 20 + 20*int(math.Floor(float64(totalFreePoints)/5))

	var pointsEarned int
	switch result {
	case "win":
		if winsCount%winThreshold == 0 {
			pointsEarned = 1
		}
	case "draw":
		if drawsCount%drawThreshold == 0 {
			pointsEarned = 1
		}
	}

	const updateStatQuery = `UPDATE user_stats
		SET value = value + $1
		WHERE user_id = $2 AND stat_type_id = $3`

	if pointsEarned > 0 {
		_, err = d.db.ExecContext(ctx, updateStatQuery, pointsEarned, userID, statTypeIDFreePoints)
		if err != nil {
			return fmt.Errorf("couldn't update user free points stat: %w", err)
		}
		_, err = d.db.ExecContext(ctx, updateStatQuery, pointsEarned, userID, statTypeIDTotalFreePoints)
		if err != nil {
			return fmt.Errorf("couldn't update user total free points stat: %w", err)
		}
	}
	return nil
}
