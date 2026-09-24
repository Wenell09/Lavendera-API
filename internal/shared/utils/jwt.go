package utils

import (
	"errors"

	"github.com/Wenell09/lavendera-api/internal/shared/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, tenantID string, role string, config config.JWTConfig) (string, error) {
	if config.SecretKey == "" {
		return "", errors.New("JWT_SECRET is not configured")
	}
	claims := JWTClaims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString(
		[]byte(config.SecretKey),
	)
}
