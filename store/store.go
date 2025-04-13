package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/yogenyslav/authgo/user"
)

// AuthStore provides methods for user authorization and account management.
type AuthStore interface {
	// InsertUser creates a new record with user data.
	InsertUser(ctx context.Context, u *user.Dao) (int64, error)
}

func applyMigrations(dialect string, migrations embed.FS, db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("setup migrations: %w", err)
	}

	if err := goose.Up(db, "schema/"+dialect); err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	return nil
}
