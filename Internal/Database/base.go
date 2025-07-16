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
	db DB
}

var (
	dbInstance *Database
)

func New(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database. Error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %s", err)
	}

	log.Printf("[%s] ✅Successfully connected to database.", time.Now().Format("15:04:05"))
	dbInstance = &Database{db: db}
	return dbInstance, nil
}

func GetInstance() *Database {
	return dbInstance
}

func (database *Database) Close() error {
	return database.db.(*sql.DB).Close()
}
