package store

import (
	"context"

	"github.com/yogenyslav/authgo/model"
)

// RoleStore provides methods to manipulate with roles.
type RoleStore interface {
	Store
	// InsertOne creates new role.
	InsertOne(ctx context.Context, name string) (int64, error)
	// FindOneByID finds a role by its id.
	FindOneByID(ctx context.Context, roleID int64) (model.RoleDao, error)
	// FindOneByName finds a role by its name.
	FindOneByName(ctx context.Context, name string) (model.RoleDao, error)
	// UpdateOne updates a role.
	UpdateOne(ctx context.Context, role model.RoleDao) error
	// DeleteOne deletes a role.
	DeleteOne(ctx context.Context, roleID int64) error
	// ListAll returns a list of all existing roles.
	ListAll(ctx context.Context) ([]model.RoleDao, error)
	// ListUserRoles returns a list of all roles asigned to a certain user.
	ListUserRoles(ctx context.Context, userID int64) ([]model.RoleDao, error)
}
