package user

import (
	"context"

	"github.com/yogenyslav/authgo/model"
)

// UserStore provides methods for user authorization and account management.
type UserStore interface {
	// InsertOne creates a new record with user data.
	InsertOne(ctx context.Context, user model.UserDao) (int64, error)
	FindOne(ctx context.Context, userID int64) (model.UserDao, error)
	UpdateOne(ctx context.Context, user model.UserDao) error
	DeleteOne(ctx context.Context, userID int64) error
	ListAll(ctx context.Context) ([]model.UserDao, error)
	SetRole(ctx context.Context, userID, roleID int64) error
	// RemoveRole(ctx context.Context, userID, roleID int64) error
	ApplyMigrations() error
}
