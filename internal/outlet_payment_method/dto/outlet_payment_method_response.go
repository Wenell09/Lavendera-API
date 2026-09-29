package dto

import (
	"time"

	"github.com/google/uuid"
)

type OutletResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OutletPaymentMethodResponse struct {
	ID            uuid.UUID      `json:"id"`
	OutletID      uuid.UUID      `json:"outlet_id"`
	Outlet        OutletResponse `json:"outlet"`
	Type          string         `json:"type"`
	ProviderName  string         `json:"provider_name"`
	AccountNumber *string        `json:"account_number"`
	AccountName   string         `json:"account_name"`
	QRImageURL    *string        `json:"qr_image_url"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type OutletPaymentMethodListResponse struct {
	Data       []OutletPaymentMethodResponse
	Pagination PaginationResponse
}
