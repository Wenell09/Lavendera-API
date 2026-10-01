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

func OutletPaymentMethodsTenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, ok := appcontext.TenantIDFromContext(ctx)
		if !ok || tenantID == uuid.Nil {
			return db.Where("1 = 0")
		}
		return db.Joins("JOIN outlets ON outlets.id = outlet_payment_methods.outlet_id").
			Where("outlets.tenant_id = ?", tenantID)
	}
}

func OutletUsersTenantScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, ok := appcontext.TenantIDFromContext(ctx)
		if !ok || tenantID == uuid.Nil {
			return db.Where("1 = 0")
		}
		return db.Joins("JOIN outlets ON outlets.id = outlet_users.outlet_id").
			Where("outlets.tenant_id = ?", tenantID)
	}
}
