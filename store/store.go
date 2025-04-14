package store

import (
	"context"
)

// Store is a base store interface.
type Store interface {
	// ApplyMigrations applies the required migrations for store.
	ApplyMigrations() error
	// StartTx starts a new transaction and puts it into context.Context.
	StartTx(ctx context.Context) (context.Context, error)
	// CommitTx commits current transaction from context.
	CommitTx(ctx context.Context) error
	// RollbackTx rolls back current transaction from context.
	RollbackTx(ctx context.Context) error
}
