package middleware

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := appcontext.RoleFromContext(c.Request.Context())
		if !ok || role == "" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				response.NewResponseError(
					http.StatusForbidden,
					"Forbidden",
					"User role not found",
				),
			)
			return
		}
		for _, allowedRole := range roles {
			if role == allowedRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			response.NewResponseError(
				http.StatusForbidden,
				"Forbidden",
				"You do not have access to this feature",
			),
		)
	}
}
