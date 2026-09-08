package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceCategory struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Tenant      Tenant    `gorm:"foreignKey:TenantID"`
	Services    []Service `gorm:"foreignKey:CategoryID"`
}
