package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed schema/postgres/*.sql
var PgMigrations embed.FS

// ApplyMigrations inits required database schema for store.
func ApplyMigrations(dialect string, migrations embed.FS, db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("setup migrations: %w", err)
	}

	if err := goose.Up(db, "schema/"+dialect); err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	return nil
}
