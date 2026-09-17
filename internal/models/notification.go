package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   *uuid.UUID
	UserID    *uuid.UUID
	Type      string
	Message   string
	IsRead    bool
	SentAt    *time.Time
	CreatedAt time.Time
	Order     *Order
	User      *User
}
