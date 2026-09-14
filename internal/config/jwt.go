package config

type JWTConfig struct {
	SecretKey string
}

func NewJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey: ENV.JWTSecret,
	}
}
