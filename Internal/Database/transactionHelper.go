package database

import (
	"context"
	"database/sql"
	"fmt"
)

func (d *Database) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := d.db.(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (d *Database) WithTransactionResult(ctx context.Context, fn func(tx *sql.Tx) (any, error)) (any, error) {
	tx, err := d.db.(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	result, err := fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("tx err: %v, rb err: %w", err, rbErr)
		}
		return nil, err
	}
	return result, tx.Commit()
}
