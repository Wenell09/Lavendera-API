package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	OutletID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	CustomerID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	DiscountID     *uuid.UUID `gorm:"type:uuid"`
	CreatedBy      *uuid.UUID `gorm:"type:uuid"`
	OrderNumber    string     `gorm:"type:varchar(100);not null"`
	TrackingToken  string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status         string     `gorm:"type:varchar(50);not null;default:'MENUNGGU_KONFIRMASI'"`
	Subtotal       int64      `gorm:"not null;default:0"`
	DiscountAmount int64      `gorm:"not null;default:0"`
	Total          int64      `gorm:"not null;default:0"`
	PaidAmount     int64      `gorm:"not null;default:0"`
	Notes          *string
	EstimatedDone  *time.Time
	CompletedAt    *time.Time
	CanceledAt     *time.Time
	CreatedAt      time.Time   `gorm:"autoCreateTime"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime"`
	Tenant         Tenant      `gorm:"foreignKey:TenantID"`
	Outlet         Outlet      `gorm:"foreignKey:OutletID"`
	Customer       Customer    `gorm:"foreignKey:CustomerID"`
	Discount       *Discount   `gorm:"foreignKey:DiscountID"`
	Creator        *User       `gorm:"foreignKey:CreatedBy"`
	Items          []OrderItem `gorm:"foreignKey:OrderID"`
	Payments       []Payment   `gorm:"foreignKey:OrderID"`
}
