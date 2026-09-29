package dto

import (
	"github.com/google/uuid"
)

type OutletPaymentMethodRequest struct {
	OutletID      uuid.UUID `json:"outlet_id" validate:"required"`
	Type          string    `json:"type" validate:"required,oneof=qris bank_transfer"`
	ProviderName  string    `json:"provider_name" validate:"required"`
	AccountNumber *string   `json:"account_number"`
	AccountName   string    `json:"account_name" validate:"required"`
	QRImageURL    *string   `json:"qr_image_url"`
	IsActive      *bool     `json:"is_active"`
}

type UpdateOutletPaymentMethodRequest struct {
	Type          *string `json:"type" validate:"omitempty,oneof=qris bank_transfer"`
	ProviderName  *string `json:"provider_name"`
	AccountNumber *string `json:"account_number"`
	AccountName   *string `json:"account_name"`
	QRImageURL    *string `json:"qr_image_url"`
	IsActive      *bool   `json:"is_active"`
}
