package store

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/yogenyslav/authgo/user"
)

// PostgresConfig holds configuration values for opening a postgres connection.
type PostgresConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	DB             string `yaml:"db"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Ssl            bool   `yaml:"ssl"`
	ConnectTimeout int    `yaml:"connect_timeout"`
}

// ConnString assembles config values into a conn string.
func (cfg *PostgresConfig) ConnString() string {
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

type postgresStore struct {
	pool *pgxpool.Pool
	cfg  PostgresConfig
}

//go:embed schema/postgres/*.sql
var pgMigrations embed.FS

// NewPostgresStore creates an instance of authStore over postgres connection.
func NewPostgresStore(cfg PostgresConfig) (*postgresStore, error) {
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

	// apply migrations
	if err = applyMigrations("postgres", pgMigrations, stdlib.OpenDBFromPool(pool)); err != nil {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return &postgresStore{
		cfg:  cfg,
		pool: pool,
	}, nil
}

const pgInsertUser = `
	insert into authgo.user(email, hash_password, username, first_name, last_name, middle_name)
	values ($1, $2, $3, $4, $5, $6)
	returning id;
`

func (s *postgresStore) InsertUser(ctx context.Context, u *user.User) (int64, error) {
	var userID int64

	err := s.pool.QueryRow(
		ctx,
		pgInsertUser,
		u.Email,
		u.HashPassword,
		u.Username,
		u.FirstName,
		u.LastName,
		u.MiddleName,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}
