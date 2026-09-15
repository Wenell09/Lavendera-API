package dto

type RegisterRequest struct {
	TenantName string `json:"tenant_name" validate:"required,min=3,max=100"`
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Email      string `json:"email" validate:"required,email,max=255"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
