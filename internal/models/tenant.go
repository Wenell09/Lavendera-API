package models

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Slug      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
