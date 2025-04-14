package authgo

import (
	"context"

	"github.com/yogenyslav/authgo/model"
)

// AuthController provides methods for user authorization.
type AuthController interface {
	// Login executes user login operation.
	Login(ctx context.Context, req model.UserLogin) (model.AuthResp, error)
	// Register executes user register operation.
	Register(ctx context.Context, req model.UserRegister) (model.AuthResp, error)
}

// UserController provides methods for manipulating with user data.
type UserController interface {
	// Me returns current user.
	Me(ctx context.Context, userID int64) (model.UserDto, error)
	// Update updates user data.
	Update(ctx context.Context, user model.UserDto) error
	// Delete deletes user.
	Delete(ctx context.Context, userID int64) error
	// ListAllUsers returns list of all existing users.
	ListAllUsers(ctx context.Context) ([]model.UserDto, error)
}

// RoleController provides methods for manipulating with user roles.
type RoleController interface {
	// SetRole assigns role to user.
	SetRole(ctx context.Context, userID, roleID int64) error
	// RemoveRole removes role from user.
	RemoveRole(ctx context.Context, userID, roleID int64) error
	// ListRoles returns list of all existing roles.
	ListRoles(ctx context.Context) ([]model.RoleDto, error)
}

// Middleware provides methods that can be used during requests to authenticate users and validate access.
type Middleware interface {
	// RequireAuth requires to pass access token with every request.
	RequireAuth(authHeader string) (model.AuthMeta, error)
	// RequireRole requires to have certain role to get access to the resource.
	RequireRole(meta model.AuthMeta, requiredRole string) error
}
