package dto

import (
	"time"

	"github.com/google/uuid"
)

type OutletUserUserResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type OutletUserOutletResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OutletUserResponse struct {
	OutletID  uuid.UUID                 `json:"outlet_id"`
	UserID    uuid.UUID                 `json:"user_id"`
	CreatedAt time.Time                 `json:"created_at"`
	Outlet    *OutletUserOutletResponse `json:"outlet,omitempty"`
	User      *OutletUserUserResponse   `json:"user,omitempty"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type OutletUserListResponse struct {
	Data       []OutletUserResponse
	Pagination PaginationResponse
}
