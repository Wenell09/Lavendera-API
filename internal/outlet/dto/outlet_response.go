package dto

import (
	"time"

	"github.com/google/uuid"
)

type OutletResponse struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Slug                 string    `json:"slug"`
	Phone                string    `json:"phone"`
	Address              string    `json:"address"`
	IsPublicOrderEnabled bool      `json:"is_public_order_enabled"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type OutletListResponse struct {
	Data       []OutletResponse
	Pagination PaginationResponse
}
