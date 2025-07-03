package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func New(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database. Error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %s", err)
	}

	log.Printf("[%s] ✅Successfully connected to database.", time.Now().Format("15:04:05"))
	return &Database{db: db}, nil
}

func (database *Database) Close() error {
	return database.db.Close()
}
