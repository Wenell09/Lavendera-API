package dto

import "github.com/google/uuid"

type CreateOutletUserRequest struct {
	OutletID uuid.UUID `json:"outlet_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}
