package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/service_category/dto"
	"github.com/google/uuid"
)

type ServiceCategoryService interface {
	FindAll(ctx context.Context) ([]dto.ServiceCategoryResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.ServiceCategoryResponse, error)
	Create(ctx context.Context, req dto.CreateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
