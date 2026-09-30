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
	Outlet    OutletUserOutletResponse `json:"outlet"`
	User      OutletUserUserResponse   `json:"user"`
	CreatedAt time.Time                `json:"created_at"`
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
