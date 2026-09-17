package models

import (
	"time"

	"github.com/google/uuid"
)

type Discount struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string
	Type      string
	Value     int64
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Tenant    Tenant
	Orders    []Order
}
