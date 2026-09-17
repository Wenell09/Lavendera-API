package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OutletID        uuid.UUID `gorm:"type:uuid;not null;index"`
	OrderID         uuid.UUID `gorm:"type:uuid;not null;index"`
	ReceivedBy      *uuid.UUID
	Amount          int64
	PaymentMethod   string
	Status          string
	ReferenceNumber *string
	Notes           *string
	PaidAt          *time.Time
	VerifiedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Outlet          Outlet
	Order           Order
	Receiver        *User `gorm:"foreignKey:ReceivedBy"`
}
