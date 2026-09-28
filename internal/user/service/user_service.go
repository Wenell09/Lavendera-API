package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/user/dto"
	"github.com/google/uuid"
)

type UserService interface {
	FindAll(ctx context.Context, filter dto.UserFilter) (*dto.UserListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
