package dto

import "github.com/google/uuid"

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ServiceResponse struct {
	ID           uuid.UUID `json:"id"`
	Category     CategoryResponse
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	MinQuantity  int64  `json:"min_quantity"`
	Unit         string `json:"unit"`
	DurationDays int64  `json:"duration_days"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ServiceListResponse struct {
	Data       []ServiceResponse `json:"data"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
}
