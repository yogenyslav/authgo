package store

import (
	"context"
)

type Store interface {
	ApplyMigrations() error
	StartTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}
