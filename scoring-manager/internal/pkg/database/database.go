package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func (db *Database) InitTransaction(ctx context.Context, nameTx string) (pgx.Tx, error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.InitTransaction: %w", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, fmt.Errorf("db.InitTransaction: %w", err)
	}

	return &txWithRelease{
		Tx:   tx,
		conn: conn,
	}, nil
}

type txWithRelease struct {
	pgx.Tx
	conn *pgxpool.Conn
}

func (tx *txWithRelease) Commit(ctx context.Context) error {
	defer tx.conn.Release()
	return tx.Tx.Commit(ctx)
}

func (tx *txWithRelease) Rollback(ctx context.Context) error {
	defer tx.conn.Release()
	return tx.Tx.Rollback(ctx)
}
