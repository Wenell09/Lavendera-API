package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID     *uuid.UUID `gorm:"type:uuid"`
	Action     string     `gorm:"type:varchar(100);not null"`
	EntityType string     `gorm:"type:varchar(100);not null"`
	OldData    []byte     `gorm:"type:jsonb"`
	NewData    []byte     `gorm:"type:jsonb"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`

	Tenant Tenant `gorm:"foreignKey:TenantID"`
	User   *User  `gorm:"foreignKey:UserID"`
}
