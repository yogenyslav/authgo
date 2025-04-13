package authgo

type AuthConfig struct {
	Jwt JwtConfig `yaml:"jwt"`
}

type JwtConfig struct {
	Secret     string `yaml:"secret"`
	Expire     int    `yaml:"expire"`
	Encryption string `yaml:"encryption"`
}
