package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey uint8

const (
	txKey contextKey = iota
)

var (
	ErrNoTxFound = errors.New("no transaction in context")
)

type postgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(pool *pgxpool.Pool) *postgresDB {
	return &postgresDB{
		pool: pool,
	}
}

func (db *postgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *postgresDB) GetConn(ctx context.Context) (*pgx.Conn, error) {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		conn, err := db.pool.Acquire(ctx)
		if err != nil {
			return nil, fmt.Errorf("acquire conn from pool: %w", err)
		}
		return conn.Conn(), nil
	}
	return tx.Conn(), nil
}

func (db *postgresDB) StartTx(ctx context.Context) (context.Context, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return ctx, fmt.Errorf("start transaction: %w", err)
	}

	return context.WithValue(ctx, txKey, tx), nil
}

func (db *postgresDB) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("commit transaction: %w", ErrNoTxFound)
	}

	return tx.Commit(ctx)
}

func (db *postgresDB) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("rollback transaction: %w", ErrNoTxFound)
	}

	err := tx.Rollback(ctx)
	switch {
	case errors.Is(err, pgx.ErrTxClosed):
		return nil
	case err != nil:
		return fmt.Errorf("rollback transaction: %w", err)
	default:
		return nil
	}
}
