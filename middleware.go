package authgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yogenyslav/authgo/model"
)

var (
	ErrMissingJwt = errors.New("missing jwt token")
	ErrForbidden  = errors.New("access to the resource is forbidden")
)

type middleware struct {
	jwt *jwtProvider
}

func NewAuthMiddleware(cfg JwtConfig) *middleware {
	return &middleware{
		jwt: newJwtProvider(cfg),
	}
}

func (m *middleware) RequireAuth(authHeader string) (model.AuthMeta, error) {
	var meta model.AuthMeta

	rawToken := strings.Split(authHeader, " ")
	if len(rawToken) != 2 {
		return meta, fmt.Errorf("parse authorization header: %w", ErrMissingJwt)
	}

	accessToken, err := m.jwt.parseAccessToken(rawToken[1])
	if err != nil {
		return meta, fmt.Errorf("invalid token: %w", err)
	}

	claims := accessToken.Claims.(jwt.MapClaims)
	rawClaims, err := json.Marshal(claims)
	if err != nil {
		return meta, fmt.Errorf("marshal token claims: %w", err)
	}

	if err := json.Unmarshal(rawClaims, &meta); err != nil {
		return meta, fmt.Errorf("unmarshal token claims: %w", err)
	}

	return meta, nil
}

func (m *middleware) RequireRole(meta model.AuthMeta, requiredRole string) error {
	for _, role := range meta.Roles {
		if role.Name == requiredRole {
			return nil
		}
	}

	return ErrForbidden
}
