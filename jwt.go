package authgo

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yogenyslav/authgo/model"
)

const (
	typeBearerToken string = "Bearer"
)

var (
	ErrJwtSignMethod = errors.New("unexpected signing method")
)

type jwtProvider struct {
	secret     string
	expire     int
	encryption []byte
}

func newJwtProvider(cfg JwtConfig) *jwtProvider {
	var encryption []byte
	if cfg.Encryption != "" {
		encryption = []byte(cfg.Encryption)
	}

	return &jwtProvider{
		secret:     cfg.Secret,
		expire:     cfg.Expire,
		encryption: encryption,
	}
}

func (j *jwtProvider) createAccessToken(meta model.AuthMeta) (string, error) {
	key := []byte(j.secret)

	jwtClaims := jwt.MapClaims{
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expire))),
		"sub":   strconv.FormatInt(meta.UserID, 10),
		"roles": meta.Roles,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := accessToken.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	if j.encryption != nil {
		return encrypt(signedToken, j.encryption)
	}

	return signedToken, nil
}

func (j *jwtProvider) parseAccessToken(accessTokenString string) (*jwt.Token, error) {
	var err error

	if j.encryption != nil {
		accessTokenString, err = decrypt(accessTokenString, j.encryption)
		if err != nil {
			return nil, fmt.Errorf("decrypt token: %w", err)
		}
	}

	accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJwtSignMethod
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
