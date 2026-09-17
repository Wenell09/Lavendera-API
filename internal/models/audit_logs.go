package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     *uuid.UUID 
	Action     string
	EntityType string
	OldData    []byte
	NewData    []byte
	CreatedAt  time.Time
	Tenant     Tenant
	User       *User
}
