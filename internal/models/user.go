package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Email        string    `gorm:"type:varchar(255);not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         string    `gorm:"type:varchar(50);not null;default:'STAFF'"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	Tenant   Tenant    `gorm:"foreignKey:TenantID"`
	Outlets  []Outlet  `gorm:"many2many:outlet_users;"`
	Orders   []Order   `gorm:"foreignKey:CreatedBy"`
	Payments []Payment `gorm:"foreignKey:ReceivedBy"`
}
