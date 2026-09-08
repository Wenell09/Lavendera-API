package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusHistory struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	Status    string     `gorm:"type:varchar(50);not null"`
	ChangedBy *uuid.UUID `gorm:"type:uuid"`
	Notes     *string    `gorm:"type:text"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	Order     Order      `gorm:"foreignKey:OrderID"`
	User      *User      `gorm:"foreignKey:ChangedBy"`
}
