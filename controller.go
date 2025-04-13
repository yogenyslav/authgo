package authgo

import (
	"context"
	"fmt"

	"github.com/yogenyslav/authgo/model"
	"github.com/yogenyslav/authgo/store/role"
	"github.com/yogenyslav/authgo/store/user"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	user user.UserStore
	role role.RoleStore
}

func NewAuthController(u user.UserStore, r role.RoleStore) (*Controller, error) {
	if err := r.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("role schema: %w", err)
	}

	if err := u.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("user schema: %w", err)
	}

	return &Controller{
		user: u,
		role: r,
	}, nil
}

func (s *Controller) Login(ctx context.Context) {

}

func (s *Controller) Register(ctx context.Context, req model.UserRegister) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("generate password hash: %w", err)
	}

	user := model.UserDao{
		Email:        req.Email,
		HashPassword: string(hashedPassword),
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
	}
	return s.user.InsertOne(ctx, user)
}

func (s *Controller) ValidateToken(ctx context.Context) {

}
