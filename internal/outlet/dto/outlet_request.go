package dto

type CreateOutletRequest struct {
	Name                 string `json:"name" validate:"required,min=5,max=255"`
	Phone                string `json:"phone" validate:"required,id_phone"`
	Address              string `json:"address" validate:"required"`
	IsPublicOrderEnabled bool   `json:"is_public_order_enabled"`
	IsActive             bool   `json:"is_active"`
}

type UpdateOutletRequest struct {
	Name                 *string `json:"name" validate:"omitempty,min=5,max=255"`
	Slug                 *string `json:"slug" validate:"omitempty,min=5,max=255"`
	Phone                *string `json:"phone" validate:"omitempty,id_phone"`
	Address              *string `json:"address" validate:"omitempty"`
	IsPublicOrderEnabled *bool   `json:"is_public_order_enabled" validate:"omitempty"`
	IsActive             *bool   `json:"is_active" validate:"omitempty"`
}
