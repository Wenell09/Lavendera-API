package models

import (
	"time"

	"github.com/google/uuid"
)

type Outlet struct {
	ID                   uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID             uuid.UUID             `gorm:"type:uuid;not null;index"`
	Name                 string                `gorm:"type:varchar(255);not null"`
	Slug                 string                `gorm:"type:varchar(100);not null"`
	Phone                *string               `gorm:"type:varchar(50)"`
	Address              *string               `gorm:"type:text"`
	IsPublicOrderEnabled bool                  `gorm:"not null;default:true"`
	IsActive             bool                  `gorm:"not null;default:true"`
	CreatedAt            time.Time             `gorm:"autoCreateTime"`
	UpdatedAt            time.Time             `gorm:"autoUpdateTime"`
	Tenant               Tenant                `gorm:"foreignKey:TenantID"`
	Users                []User                `gorm:"many2many:outlet_users;"`
	Orders               []Order               `gorm:"foreignKey:OutletID"`
	Payments             []Payment             `gorm:"foreignKey:OutletID"`
	PaymentMethods       []OutletPaymentMethod `gorm:"foreignKey:OutletID"`
}
