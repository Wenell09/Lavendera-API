package middleware

import (
	"net/http"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/config"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTMiddleware(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.NewResponseError(
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
				response.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"Format token salah (gunakan Bearer <token>)",
				),
			)
			return
		}
		tokenString := parts[1]
		claims := &utils.JWTClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				// Hanya izinkan algoritma HS256.
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtConfig.SecretKey), nil
			},
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.NewResponseError(
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
				response.NewResponseError(
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
				response.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"tenant_id tidak ditemukan dalam token",
				),
			)
			return
		}
		// Simpan identity ke request context.
		rawUserID := claims.UserID
		userID, err := uuid.Parse(rawUserID)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"Invalid User ID format",
				),
			)
			return
		}
		rawTenantID := claims.TenantID
		tenantID, err := uuid.Parse(rawTenantID)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"Invalid Tenant ID format",
				),
			)
			return
		}
		ctx := c.Request.Context()
		ctx = appcontext.WithUserID(ctx, userID)
		ctx = appcontext.WithTenantID(ctx, tenantID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
