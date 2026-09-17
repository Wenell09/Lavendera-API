package models

import (
	"time"

	"github.com/google/uuid"
)

type Outlet struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID             uuid.UUID `gorm:"type:uuid;not null;index"`
	Name                 string
	Slug                 string
	Phone                *string
	Address              *string
	IsPublicOrderEnabled bool
	IsActive             bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Tenant               Tenant
	Users                []User
	Orders               []Order
	Payments             []Payment
	PaymentMethods       []OutletPaymentMethod
}
