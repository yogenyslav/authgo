package store

import (
	"context"

	"github.com/yogenyslav/authgo/model"
)

// UserStore provides methods for user authorization and account management.
type UserStore interface {
	Store
	// InsertOne creates a new record with user data.
	InsertOne(ctx context.Context, user model.UserDao) (int64, error)
	// FindOneByID finds a user by its id.
	FindOneByID(ctx context.Context, id int64) (model.UserDao, error)
	// FindOneByEmail finds a user by its email.
	FindOneByEmail(ctx context.Context, email string) (model.UserDao, error)
	// UpdateOne updates user.
	UpdateOne(ctx context.Context, user model.UserDao) error
	// DeleteOne deletes user.
	DeleteOne(ctx context.Context, userID int64) error
	// ListAll return the list of all existing users.
	ListAll(ctx context.Context) ([]model.UserDao, error)
	// SetRole assigns role to user.
	SetRole(ctx context.Context, userID, roleID int64) error
	// RemoveRole removes role from user.
	RemoveRole(ctx context.Context, userID, roleID int64) error
}
