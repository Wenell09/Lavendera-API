package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
