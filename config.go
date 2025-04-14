package authgo

import "github.com/yogenyslav/authgo/store/postgres"

// AuthConfig is a top-level config for authgo package that holds other nested configs.
type AuthConfig struct {
	Jwt      JwtConfig       `yaml:"jwt"`
	Postgres postgres.Config `yaml:"postgres"`
}

// JwtConfig is a config for jwt module.
type JwtConfig struct {
	Secret     string `yaml:"secret"`
	Expire     int    `yaml:"expire"`
	Encryption string `yaml:"encryption"`
}
