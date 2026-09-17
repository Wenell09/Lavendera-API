package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ServiceID uuid.UUID
	Price     int64
	Quantity  int64
	Subtotal  int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Order     Order
	Service   Service
}
