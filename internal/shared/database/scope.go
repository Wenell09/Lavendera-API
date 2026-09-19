package database

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, ok := appcontext.TenantIDFromContext(ctx)
		if !ok || tenantID == uuid.Nil {
			return db.Where("1 = 0")
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}
