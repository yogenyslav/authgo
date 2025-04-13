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
	FindOneByEmail(ctx context.Context, email string) (model.UserDao, error)
	UpdateOne(ctx context.Context, user model.UserDao) error
	DeleteOne(ctx context.Context, userID int64) error
	ListAll(ctx context.Context) ([]model.UserDao, error)
	SetRole(ctx context.Context, userID, roleID int64) error
	// RemoveRole(ctx context.Context, userID, roleID int64) error
}
