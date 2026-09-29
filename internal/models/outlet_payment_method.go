package models

import (
	"time"

	"github.com/google/uuid"
)

type OutletPaymentMethod struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID      uuid.UUID `gorm:"type:uuid;not null;index"`
	OutletID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Type          string
	ProviderName  string
	AccountNumber *string
	AccountName   string
	QRImageURL    *string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Outlet        Outlet
}
