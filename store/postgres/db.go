package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey uint8

const (
	txKey contextKey = iota
)

var (
	// ErrNoTxFound is an error when requested a transaction mode, but Tx is not found in context.Context.
	ErrNoTxFound = errors.New("no transaction in context")
)

// Config holds configuration values for opening a postgres connection.
type Config struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	DB             string `yaml:"db"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Ssl            bool   `yaml:"ssl"`
	ConnectTimeout int    `yaml:"connect_timeout"`
}

// ConnString assembles config values into a conn string.
func (cfg *Config) ConnString() string {
	sslMode := "disable"
	if cfg.Ssl {
		sslMode = "enable"
	}
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		sslMode,
	)
}

type postgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates new postgresDB (wrapper for pgxpool.Pool).
func NewPostgresDB(cfg Config) (*postgresDB, error) {
	// init postgres config from pgx
	pgConfig, err := pgxpool.ParseConfig(cfg.ConnString())
	if err != nil {
		return nil, fmt.Errorf("parse postgres connection string: %w", err)
	}

	// set few connection options
	pgConfig.ConnConfig.ConnectTimeout = time.Duration(cfg.ConnectTimeout)
	pgConfig.ConnConfig.ValidateConnect = func(ctx context.Context, conn *pgconn.PgConn) error {
		return conn.Ping(ctx)
	}

	// create pgx pool for postgres
	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	return &postgresDB{
		pool: pool,
	}, nil
}

// GetPool returns the underlying *pgxpool.Pool.
func (db *postgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}

// GetConn returns either a Tx (if has one), or acquires a new *pgx.Conn from *pgxpool.Pool.
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

// StartTx starts a new transaction and puts into the context.Context.
func (db *postgresDB) StartTx(ctx context.Context) (context.Context, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return ctx, fmt.Errorf("start transaction: %w", err)
	}

	return context.WithValue(ctx, txKey, tx), nil
}

// CommitTx commits current transaction from context.
func (db *postgresDB) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("commit transaction: %w", ErrNoTxFound)
	}

	return tx.Commit(ctx)
}

// RollbackTx rolls back current transaction from context.
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
