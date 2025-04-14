package authgo

import (
	"context"
	"errors"
	"fmt"

	"github.com/yogenyslav/authgo/model"
	"github.com/yogenyslav/authgo/store"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

// controller provides methods to manipulate with user and its roles.
type controller struct {
	cfg  AuthConfig
	user store.UserStore
	role store.RoleStore
	jwt  *jwtProvider
}

// NewAuthController is a constructor for Controller.
func NewAuthController(cfg AuthConfig, u store.UserStore, r store.RoleStore) (*controller, error) {
	if err := r.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("role schema: %w", err)
	}

	if err := u.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("user schema: %w", err)
	}

	jwt := newJwtProvider(cfg.Jwt)

	return &controller{
		cfg:  cfg,
		user: u,
		role: r,
		jwt:  jwt,
	}, nil
}

func (ctrl *controller) Login(ctx context.Context, req model.UserLogin) (model.AuthResp, error) {
	var resp model.AuthResp

	user, err := ctrl.user.FindOneByEmail(ctx, req.Email)
	if err != nil {
		return resp, fmt.Errorf("find user: %w", err)
	}

	if !verifyPassword(user.HashPassword, req.Password) {
		return resp, fmt.Errorf("verify password: %w", err)
	}

	rolesDB, err := ctrl.role.ListUserRoles(ctx, user.ID)
	if err != nil {
		return resp, fmt.Errorf("list user roles: %w", err)
	}

	roles := make([]model.RoleDto, 0, len(rolesDB))
	for _, role := range rolesDB {
		roles = append(roles, role.ToDto())
	}

	meta := model.AuthMeta{
		UserID: user.ID,
		Roles:  roles,
	}
	accessToken, err := ctrl.jwt.createAccessToken(meta)
	if err != nil {
		return resp, fmt.Errorf("create access token: %w", err)
	}

	resp.Meta = meta
	resp.Token = accessToken
	resp.Type = typeBearerToken

	return resp, nil
}

func (ctrl *controller) Register(ctx context.Context, req model.UserRegister) (model.AuthResp, error) {
	var resp model.AuthResp

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return resp, fmt.Errorf("hash password: %w", err)
	}

	ctx, err = ctrl.user.StartTx(ctx)
	if err != nil {
		return resp, fmt.Errorf("user store transaction: %w", err)
	}
	defer func() {
		if err := ctrl.user.RollbackTx(ctx); err != nil {
			panic(fmt.Errorf("rollback user store transaction: %w", err))
		}
	}()

	user := model.UserDao{
		Email:        req.Email,
		HashPassword: hashedPassword,
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
	}
	userID, err := ctrl.user.InsertOne(ctx, user)
	if err != nil {
		return resp, fmt.Errorf("insert user: %w", err)
	}

	role, err := ctrl.role.FindOneByName(ctx, model.DefaultRole)
	if err != nil {
		return resp, fmt.Errorf("find role by name: %w", err)
	}

	if err = ctrl.user.SetRole(ctx, userID, role.ID); err != nil {
		return resp, fmt.Errorf("set role: %w", err)
	}

	meta := model.AuthMeta{
		UserID: userID,
		Roles: []model.RoleDto{{
			ID:   userID,
			Name: model.DefaultRole,
		}},
	}
	accessToken, err := ctrl.jwt.createAccessToken(meta)
	if err != nil {
		return resp, fmt.Errorf("create access token: %w", err)
	}

	if err = ctrl.user.CommitTx(ctx); err != nil {
		return resp, fmt.Errorf("commit user transaction: %w", err)
	}

	resp.Meta = meta
	resp.Token = accessToken
	resp.Type = typeBearerToken

	return resp, nil
}

func (ctrl *controller) Me(ctx context.Context, userID int64) (model.UserDto, error) {
	userDB, err := ctrl.user.FindOneByID(ctx, userID)
	if err != nil {
		return model.UserDto{}, fmt.Errorf("current user: %w", err)
	}

	return userDB.ToDto(), nil
}

func (ctrl *controller) Update(ctx context.Context, u model.UserDto) error {
	user := model.UserDao{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		MiddleName: u.MiddleName,
	}
	return ctrl.user.UpdateOne(ctx, user)
}

func (ctrl *controller) Delete(ctx context.Context, userID int64) error {
	return ctrl.user.DeleteOne(ctx, userID)
}

func (ctrl *controller) ListAllUsers(ctx context.Context) ([]model.UserDto, error) {
	usersDB, err := ctrl.user.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all: %w", err)
	}

	users := make([]model.UserDto, 0, len(usersDB))
	for _, user := range usersDB {
		users = append(users, user.ToDto())
	}
	return users, nil
}

func (ctrl *controller) SetRole(ctx context.Context, userID, roleID int64) error {
	return ctrl.user.SetRole(ctx, userID, roleID)
}

func (ctrl *controller) RemoveRole(ctx context.Context, userID, roleID int64) error {
	return ctrl.user.RemoveRole(ctx, userID, roleID)
}

func (ctrl *controller) ListRoles(ctx context.Context) ([]model.RoleDto, error) {
	rolesDB, err := ctrl.role.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all roles: %w", err)
	}

	roles := make([]model.RoleDto, 0, len(rolesDB))
	for _, role := range rolesDB {
		roles = append(roles, role.ToDto())
	}
	return roles, nil
}
