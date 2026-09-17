package database

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/shared/middleware"
	"gorm.io/gorm"
)

func TenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, ok := middleware.TenantIDFromContext(ctx)
		if !ok || tenantID == "" {
			return db.Where("1 = 0")
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}
