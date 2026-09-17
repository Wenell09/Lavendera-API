package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string
	Phone     string
	Address   *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tenant    Tenant
	Orders    []Order
}
