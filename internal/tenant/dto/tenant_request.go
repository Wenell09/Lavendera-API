package dto

type UpdateTenantRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=100"`
	Slug *string `json:"slug" validate:"omitempty,min=3,max=100"`
}
