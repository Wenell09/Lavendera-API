package models

import (
	"time"

	"github.com/google/uuid"
)

type Discount struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Type      string    `gorm:"type:varchar(20);not null"`
	Value     int64     `gorm:"not null;default:0"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID"`
	Orders    []Order   `gorm:"foreignKey:DiscountID"`
}
