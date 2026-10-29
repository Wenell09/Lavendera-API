package dto

import (
	"time"

	"github.com/google/uuid"
)

type OutletUserUserResponse struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Role     string     `json:"role"`
	IsActive bool       `json:"is_active"`
	AssignAt *time.Time `json:"assigned_at,omitempty"`
}

type OutletUserOutletResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type OutletUsersResponse struct {
	Outlet OutletUserOutletResponse `json:"outlet"`
	Users  []OutletUserUserResponse `json:"users"`
}

type UserOutletsResponse struct {
	User   OutletUserUserResponse   `json:"user"`
	Outlet OutletUserOutletResponse `json:"outlet"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type OutletUserListResponse struct {
	Data       OutletUsersResponse
	Pagination PaginationResponse
}
