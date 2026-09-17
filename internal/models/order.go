package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index"`
	OutletID       uuid.UUID `gorm:"type:uuid;not null;index"`
	CustomerID     uuid.UUID `gorm:"type:uuid;not null;index"`
	DiscountID     *uuid.UUID
	CreatedBy      *uuid.UUID
	OrderNumber    string
	TrackingToken  string
	Status         string
	Subtotal       int64
	DiscountAmount int64
	Total          int64
	PaidAmount     int64
	Notes          *string
	EstimatedDone  *time.Time
	CompletedAt    *time.Time
	CanceledAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Tenant         Tenant
	Outlet         Outlet
	Customer       Customer
	Discount       *Discount
	Creator        *User `gorm:"foreignKey:CreatedBy"`
	Items          []OrderItem
	Payments       []Payment
}
