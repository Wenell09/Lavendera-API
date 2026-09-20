package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/google/uuid"
)

type Service interface {
	FindAll(ctx context.Context, filter dto.ServiceFilter) (*dto.ServiceListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.ServiceResponse, error)
	Create(ctx context.Context, request dto.CreateServiceRequest) (*dto.ServiceResponse, error)
	Update(ctx context.Context, id uuid.UUID, request dto.UpdateServiceRequest) (*dto.ServiceResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
