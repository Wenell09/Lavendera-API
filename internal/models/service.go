package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID    uuid.UUID        `gorm:"type:uuid;not null;index"`
	CategoryID  *uuid.UUID       `gorm:"type:uuid"`
	Name        string           `gorm:"type:varchar(255);not null"`
	Price       int64            `gorm:"not null;default:0"`
	MinQuantity int64            `gorm:"not null;default:1000"`
	Unit        string           `gorm:"type:varchar(20);not null;default:'gram'"`
	Description *string          `gorm:"type:text"`
	IsActive    bool             `gorm:"not null;default:true"`
	CreatedAt   time.Time        `gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime"`
	Tenant      Tenant           `gorm:"foreignKey:TenantID"`
	Category    *ServiceCategory `gorm:"foreignKey:CategoryID"`
	Items       []OrderItem      `gorm:"foreignKey:ServiceID"`
}
