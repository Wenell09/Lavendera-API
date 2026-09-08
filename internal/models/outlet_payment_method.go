package models

import (
	"time"

	"github.com/google/uuid"
)

type OutletPaymentMethod struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OutletID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Type          string    `gorm:"type:varchar(50);not null"`
	ProviderName  string    `gorm:"type:varchar(100);not null"`
	AccountNumber *string   `gorm:"type:varchar(100)"`
	AccountName   string    `gorm:"type:varchar(255);not null"`
	QRImageURL    *string   `gorm:"type:varchar(500)"`
	IsActive      bool      `gorm:"not null;default:true"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	Outlet        Outlet    `gorm:"foreignKey:OutletID"`
}
