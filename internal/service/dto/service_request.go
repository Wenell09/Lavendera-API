package dto

import "github.com/google/uuid"

type CreateServiceRequest struct {
	CategoryID   uuid.UUID `json:"category_id" validate:"required"`
	Name         string    `json:"name" validate:"required,min=3,max=255"`
	Price        int64     `json:"price" validate:"gte=0"`
	MinQuantity  int64     `json:"min_quantity" validate:"gt=0"`
	Unit         string    `json:"unit" validate:"required"`
	DurationDays int64     `json:"duration_days" validate:"gt=0"`
	IsActive     bool     `json:"is_active"`
}

type UpdateServiceRequest struct {
	CategoryID   *uuid.UUID `json:"category_id"`
	Name         *string    `json:"name" validate:"omitempty,min=3,max=255"`
	Price        *int64     `json:"price" validate:"omitempty,min=0"`
	MinQuantity  *int64     `json:"min_quantity" validate:"omitempty,min=1"`
	Unit         *string    `json:"unit" validate:"omitempty"`
	DurationDays *int64     `json:"duration_days" validate:"omitempty,min=0"`
	IsActive     *bool      `json:"is_active"`
}
