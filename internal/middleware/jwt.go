package middleware

import (
	"net/http"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/auth/service"
	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/Wenell09/lavendera-api/internal/contextutil"
	"github.com/Wenell09/lavendera-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"Header otorisasi tidak ditemukan",
				),
			)
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"Format token salah (gunakan Bearer <token>)",
				),
			)
			return
		}
		tokenString := parts[1]
		claims := &service.JWTClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				// Hanya izinkan HS256
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtConfig.SecretKey), nil
			},
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"Token tidak valid",
				),
			)
			return
		}
		if claims.UserID == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"user_id tidak ditemukan dalam token",
				),
			)
			return
		}
		if claims.TenantID == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"tenant_id tidak ditemukan dalam token",
				),
			)
			return
		}
		ctx := c.Request.Context()
		ctx = contextutil.WithUserID(
			ctx,
			claims.UserID,
		)
		ctx = contextutil.WithTenantID(
			ctx,
			claims.TenantID,
		)
		c.Request = c.Request.WithContext(ctx)
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Next()
	}
}
