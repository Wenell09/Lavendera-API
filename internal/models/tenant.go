package models

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string     `gorm:"type:varchar(255);not null"`
	Slug      string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	LogoURL   *string    `gorm:"type:varchar(500)"`
	Domain    *string    `gorm:"type:varchar(255);uniqueIndex"`
	Email     string     `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	Users     []User     `gorm:"foreignKey:TenantID"`
	Outlets   []Outlet   `gorm:"foreignKey:TenantID"`
	Customers []Customer `gorm:"foreignKey:TenantID"`
	Services  []Service  `gorm:"foreignKey:TenantID"`
	Discounts []Discount `gorm:"foreignKey:TenantID"`
	Orders    []Order    `gorm:"foreignKey:TenantID"`
}
