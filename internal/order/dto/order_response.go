package dto

import (
	"time"

	"github.com/google/uuid"
)

type CustomerSummary struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Address *string   `json:"address,omitempty"`
}

type OutletSummary struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Address string    `json:"address"`
}

type DiscountSummary struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Value int64     `json:"value"`
}

type UserSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ServiceCategorySummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ServiceSummary struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Price        int64                  `json:"price"`
	DurationDays int64                  `json:"duration_days"`
	Category     ServiceCategorySummary `json:"category"`
}

type OrderItemResponse struct {
	ID       uuid.UUID      `json:"id"`
	Price    int64          `json:"price"`
	Quantity float64        `json:"quantity"`
	Unit     string         `json:"unit"`
	Subtotal int64          `json:"subtotal"`
	Service  ServiceSummary `json:"service"`
}

type PaymentResponse struct {
	ID              uuid.UUID  `json:"id"`
	Amount          int64      `json:"amount"`
	PaymentMethod   string     `json:"payment_method"`
	Status          string     `json:"status"`
	ReferenceNumber *string    `json:"reference_number,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OrderResponse struct {
	ID             uuid.UUID           `json:"id"`
	TenantID       uuid.UUID           `json:"tenant_id"`
	OrderNumber    string              `json:"order_number"`
	TrackingToken  string              `json:"tracking_token"`
	Status         string              `json:"status"`
	PaymentStatus  string              `json:"payment_status"`
	Subtotal       int64               `json:"subtotal"`
	DiscountAmount int64               `json:"discount_amount"`
	TotalAmount    int64               `json:"total_amount"`
	PaidAmount     int64               `json:"paid_amount"`
	ChangeAmount   int64               `json:"change_amount"`
	Outlet         OutletSummary       `json:"outlet"`
	Customer       CustomerSummary     `json:"customer"`
	Discount       *DiscountSummary    `json:"discount,omitempty"`
	Creator        *UserSummary        `json:"creator,omitempty"`
	Items          []OrderItemResponse `json:"items,omitempty"`
	Payments       []PaymentResponse   `json:"payments"`
	Notes          *string             `json:"notes,omitempty"`
	EstimatedDone  time.Time           `json:"estimated_done"`
	CompletedAt    *time.Time          `json:"completed_at,omitempty"`
	CanceledAt     *time.Time          `json:"canceled_at,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type OrderListResponse struct {
	Data       []OrderResponse
	Pagination PaginationResponse
}
