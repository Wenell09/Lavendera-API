package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ServiceID uuid.UUID `gorm:"type:uuid;not null"`
	Price     int64     `gorm:"not null;default:0"`
	Quantity  int64     `gorm:"not null;default:0"`
	Subtotal  int64     `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Order     Order     `gorm:"foreignKey:OrderID"`
	Service   Service   `gorm:"foreignKey:ServiceID"`
}
