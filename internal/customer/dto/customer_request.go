package dto

type CreateCustomerRequest struct {
	Name    string  `json:"name" validate:"required,min=3,max=150"`
	Phone   string  `json:"phone" validate:"required,id_phone"`
	Address *string `json:"address" validate:"omitempty,max=500"`
}

type UpdateCustomerRequest struct {
	Name    *string `json:"name" validate:"omitempty,min=3,max=150"`
	Phone   *string `json:"phone" validate:"omitempty,id_phone"`
	Address *string `json:"address" validate:"omitempty,max=500"`
}
