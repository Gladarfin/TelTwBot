package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetOrCreateUser(t *testing.T) {
	t.Run("existing user", func(t *testing.T) {
		db, mock := SetupMockDB(t)
		defer db.Close()

		ctx := context.Background()
		username := "existinguser"
		userID := 1

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
		mock.ExpectCommit()

		id, err := db.getOrCreateUser(ctx, username)
		require.NoError(t, err)
		require.Equal(t, userID, id)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("new user", func(t *testing.T) {
		db, mock := SetupMockDB(t)
		defer db.Close()

		ctx := context.Background()
		username := "newuser"
		userID := 2

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery("INSERT INTO users .* RETURNING id").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
		mock.ExpectCommit()

		id, err := db.getOrCreateUser(ctx, username)
		require.NoError(t, err)
		require.Equal(t, userID, id)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetOrCreateUserStats(t *testing.T) {
	t.Run("success - new user", func(t *testing.T) {
		db, mock := SetupMockDB(t)
		defer db.Close()

		ctx := context.Background()
		username := "testuser"
		userID := 1
		now := time.Now()

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery("INSERT INTO users .* RETURNING id").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		mock.ExpectQuery("SELECT name FROM stat_types").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).
				AddRow("stat1").
				AddRow("stat2"))

		mock.ExpectExec("INSERT INTO user_stats .* ON CONFLICT .*").
			WithArgs(userID, "stat1").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("INSERT INTO user_stats .* ON CONFLICT .*").
			WithArgs(userID, "stat2").
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT s.name, us.value, us.updated_at .*").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"name", "value", "updated_at"}).
				AddRow("stat1", 10, now).
				AddRow("stat2", 20, now))

		mock.ExpectCommit()

		stats, err := db.GetOrCreateUserStats(ctx, username)
		require.NoError(t, err)
		require.Len(t, stats, 2)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database failure", func(t *testing.T) {
		db, mock := SetupMockDB(t)
		defer db.Close()

		ctx := context.Background()
		username := "testuser"

		mock.ExpectBegin().WillReturnError(errors.New("tx error"))

		_, err := db.GetOrCreateUserStats(ctx, username)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to begin transaction")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - user query fails", func(t *testing.T) {
		db, mock := SetupMockDB(t)
		defer db.Close()

		ctx := context.Background()
		username := "testuser"

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnError(errors.New("db error"))

		mock.ExpectRollback()

		_, err := db.GetOrCreateUserStats(ctx, username)
		require.Error(t, err)
		require.Contains(t, err.Error(), "user setup failed")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
