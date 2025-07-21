package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	now := time.Now()
	expectedUser := &User{
		ID:        1,
		Username:  "testuser",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`
		INSERT INTO users \(username\)
		VALUES \(\$1\)
		ON CONFLICT \(username\) DO UPDATE SET updated_at = NOW\(\)
		RETURNING id, username, created_at, updated_at
	`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.CreatedAt, expectedUser.UpdatedAt))
	mock.ExpectCommit()

	database := &Database{db: db}

	user, err := database.CreateUser(context.Background(), "testuser")

	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`
		INSERT INTO users \(username\)
		VALUES \(\$1\)
		ON CONFLICT \(username\) DO UPDATE SET updated_at = NOW\(\)
		RETURNING id, username, created_at, updated_at
	`).
		WithArgs("testuser").
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	database := &Database{db: db}

	user, err := database.CreateUser(context.Background(), "testuser")

	require.Error(t, err)
	require.Nil(t, user)
	require.EqualError(t, err, "failed to create user: database error")
	require.NoError(t, mock.ExpectationsWereMet())
}
