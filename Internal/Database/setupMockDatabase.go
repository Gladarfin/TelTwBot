package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func SetupMockDB(t *testing.T) (*Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectBegin().WillReturnError(nil)

	return &Database{db: db}, mock
}
