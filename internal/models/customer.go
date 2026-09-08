package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Phone     string    `gorm:"type:varchar(50);not null"`
	Address   *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID"`
	Orders    []Order   `gorm:"foreignKey:CustomerID"`
}
