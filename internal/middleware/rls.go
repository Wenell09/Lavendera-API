package middleware

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/contextutil"
	"github.com/Wenell09/lavendera-api/internal/database"
	"github.com/Wenell09/lavendera-api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RLSMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := contextutil.TenantIDFromContext(c.Request.Context())
		if !exists {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.NewResponseError(
					http.StatusUnauthorized,
					"Unauthorized access",
					"tenant_id not found",
				),
			)
			return
		}
		tx := db.Begin()
		if tx.Error != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.NewResponseError(
					http.StatusInternalServerError,
					"Internal Server Error",
					"failed to begin transaction",
				),
			)
			return
		}
		err := tx.Exec(
			"SELECT set_config('app.current_tenant_id', ?, true)",
			tenantID,
		).Error
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.NewResponseError(
					http.StatusInternalServerError,
					"Internal Server Error",
					"failed to set tenant context",
				),
			)
			return
		}
		ctx := database.WithTx(
			c.Request.Context(),
			tx,
		)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		if c.Writer.Status() >= 400 {
			tx.Rollback()
			return
		}
		if err := tx.Commit().Error; err != nil {
			return
		}
	}
}
