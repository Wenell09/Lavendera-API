package dto

import (
	"time"

	"github.com/google/uuid"
)

type OutletUserUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
}

type OutletUserOutletResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type OutletUserResponse struct {
	OutletID  uuid.UUID                `json:"outlet_id"`
	UserID    uuid.UUID                `json:"user_id"`
	CreatedAt time.Time                `json:"created_at"`
	Outlet    *OutletUserOutletResponse `json:"outlet,omitempty"`
	User      *OutletUserUserResponse   `json:"user,omitempty"`
}
