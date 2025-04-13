package store

import (
	"context"

	"github.com/yogenyslav/authgo/model"
)

type RoleStore interface {
	Store
	InsertOne(ctx context.Context, name string) (int64, error)
	FindOneByID(ctx context.Context, roleID int64) (model.RoleDao, error)
	FindOneByName(ctx context.Context, name string) (model.RoleDao, error)
	UpdateOne(ctx context.Context, role model.RoleDao) error
	DeleteOne(ctx context.Context, roleID int64) error
	ListAll(ctx context.Context) ([]model.RoleDao, error)
	ListUserRoles(ctx context.Context, userID int64) ([]model.RoleDao, error)
}
