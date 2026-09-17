package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index"`
	CategoryID   uuid.UUID
	Name         string
	Price        int64
	MinQuantity  int64
	Unit         string
	DurationDays int64
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Category     ServiceCategory
}
