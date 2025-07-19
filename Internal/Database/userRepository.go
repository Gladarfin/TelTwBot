package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Database) CreateUser(ctx context.Context, username string) (*User, error) {
	var user User
	err := d.WithTransaction(ctx, func(tx *sql.Tx) error {
		const query = `
			INSERT INTO users (username)
			VALUES ($1)
			ON CONFLICT (username) DO UPDATE SET updated_at = NOW()
			RETURNING id, username, created_at, updated_at
		`

		return tx.QueryRowContext(ctx, query, username).Scan(
			&user.ID,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (d *Database) GetUser(ctx context.Context, username string) (*User, error) {
	var user User
	err := d.WithTransaction(ctx, func(tx *sql.Tx) error {
		const query = `
				SELECT id, username, created_at, updated_at
				FROM users
				WHERE username = $1				
		`

		return tx.QueryRowContext(ctx, query, username).Scan(
			&user.ID,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
