package dto

type CreateDiscountRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Type     string `json:"type" validate:"required,oneof=percentage fixed"`
	Value    int64  `json:"value" validate:"required,gt=0"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type UpdateDiscountRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=3,max=255"`
	Type     *string `json:"type" validate:"omitempty,oneof=percentage fixed"`
	Value    *int64  `json:"value" validate:"omitempty,min=0"`
	IsActive *bool   `json:"is_active"`
}
