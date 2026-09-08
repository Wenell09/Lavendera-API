package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OutletID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	OrderID         uuid.UUID  `gorm:"type:uuid;not null;index"`
	ReceivedBy      *uuid.UUID `gorm:"type:uuid"`
	Amount          int64      `gorm:"not null;default:0"`
	PaymentMethod   string     `gorm:"type:varchar(50);not null"`
	Status          string     `gorm:"type:varchar(20);not null;default:'pending'"`
	ReferenceNumber *string    `gorm:"type:varchar(100)"`
	Notes           *string    `gorm:"type:text"`
	PaidAt          *time.Time
	VerifiedAt      *time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	Outlet          Outlet    `gorm:"foreignKey:OutletID"`
	Order           Order     `gorm:"foreignKey:OrderID"`
	Receiver        *User     `gorm:"foreignKey:ReceivedBy"`
}
