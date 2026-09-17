package models

import (
	"time"

	"github.com/google/uuid"
)

type OutletUser struct {
	OutletID  uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	Outlet    Outlet
	User      User
}
