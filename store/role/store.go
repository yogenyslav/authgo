package role

import (
	"context"

	"github.com/yogenyslav/authgo/model"
)

type RoleStore interface {
	InsertOne(ctx context.Context, name string) (int64, error)
	FindOne(ctx context.Context, roleID int64) (model.RoleDao, error)
	UpdateOne(ctx context.Context, role model.RoleDao) error
	DeleteOne(ctx context.Context, roleID int64) error
	ListAll(ctx context.Context) ([]model.RoleDao, error)
	ApplyMigrations() error
}
