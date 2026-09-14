package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/auth/dto"
)

type AuthService interface {
	Register(
		ctx context.Context,
		req dto.RegisterRequest,
	) (*dto.RegisterResponse, error)

	Login(
		ctx context.Context,
		req dto.LoginRequest,
	) (*dto.LoginResponse, error)
}
