package dto

import (
	"github.com/google/uuid"
)

type CreateOrderItemRequest struct {
	ServiceID uuid.UUID `json:"service_id" validate:"required"`
	Quantity  float64   `json:"quantity" validate:"required,gt=0"`
}

type CreateCashierOrderRequest struct {
	OutletID      uuid.UUID                `json:"outlet_id" validate:"required"`
	CustomerID    uuid.UUID                `json:"customer_id" validate:"required"`
	DiscountID    *uuid.UUID               `json:"discount_id" validate:"omitempty"`
	PaymentMethod string                   `json:"payment_method" validate:"required,oneof=cash bank_transfer qris"`
	PaidAmount    *int64                   `json:"paid_amount" validate:"omitempty,gte=0"`
	Notes         *string                  `json:"notes" validate:"omitempty,max=1000"`
	Status        string                   `json:"status" validate:"omitempty,oneof=MENUNGGU_KONFIRMASI DIPROSES"`
	Items         []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreatePublicOrderRequest struct {
	OutletSlug    string                   `json:"outlet_slug" validate:"required,min=1,max=100"`
	Name          string                   `json:"name" validate:"required,min=3,max=150"`
	Phone         string                   `json:"phone" validate:"required,id_phone"`
	Address       *string                  `json:"address" validate:"omitempty,max=500"`
	DiscountID    *uuid.UUID               `json:"discount_id" validate:"omitempty"`
	PaymentMethod string                   `json:"payment_method" validate:"required,oneof=cash bank_transfer qris"`
	PaidAmount    *int64                   `json:"paid_amount" validate:"omitempty,gte=0"`
	Notes         *string                  `json:"notes" validate:"omitempty,max=1000"`
	Items         []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=MENUNGGU_KONFIRMASI DIPROSES SELESAI DIAMBIL"`
}
