package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=150"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Role     string `json:"role" validate:"required,oneof=ADMIN STAFF"`
	IsActive *bool  `json:"is_active"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=150"`
	Email    *string `json:"email" validate:"omitempty,email,max=255"`
	Password *string `json:"password" validate:"omitempty,min=8,max=72"`
	Role     *string `json:"role" validate:"omitempty,oneof=ADMIN STAFF"`
	IsActive *bool   `json:"is_active" validate:"omitempty"`
}
