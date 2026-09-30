package dto

import "github.com/google/uuid"

type RegisterResponse struct {
	Tenant TenantResponse `json:"tenant"`
	User   UserResponse   `json:"user"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type TenantResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Slug  string    `json:"slug"`
	Email string    `json:"email"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
