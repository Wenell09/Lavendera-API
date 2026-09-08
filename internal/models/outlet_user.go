package models

import (
	"time"

	"github.com/google/uuid"
)

type OutletUser struct {
	OutletID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Outlet    Outlet    `gorm:"foreignKey:OutletID"`
	User      User      `gorm:"foreignKey:UserID"`
}
