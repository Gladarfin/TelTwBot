package database

import (
	"context"
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

	const query = `
			INSERT INTO users (username)
			VALUES ($1)
			ON CONFLICT (username) DO UPDATE SET updated_at = NOW()
			RETURNING id, username, created_at, updated_at
	`

	var user User
	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("can't create or update current user: %w", err)
	}

	return &user, nil
}

func (d *Database) GetUser(ctx context.Context, username string) (*User, error) {
	const query = `
			SELECT id, username, created_at, updated_at
			FROM users
			WHERE username = $1				
	`

	var user User
	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("current user doesn't exist in database. Error: %w", err)
	}

	return &user, nil
}
