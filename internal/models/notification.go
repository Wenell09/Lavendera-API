package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   *uuid.UUID `gorm:"type:uuid"`
	UserID    *uuid.UUID `gorm:"type:uuid"`
	Type      string     `gorm:"type:varchar(50);not null"`
	Message   string     `gorm:"type:text;not null"`
	SentAt    *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Order     *Order    `gorm:"foreignKey:OrderID"`
	User      *User     `gorm:"foreignKey:UserID"`
}
