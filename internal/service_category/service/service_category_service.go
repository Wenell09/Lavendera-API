package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/service_category/dto"
)

type ServiceCategoryService interface {
	FindAll(
		ctx context.Context,
	) (*dto.ServiceCategoryListResponse, error)

	FindByID(
		ctx context.Context,
		id string,
	) (*dto.ServiceCategoryResponse, error)

	Create(
		ctx context.Context,
		req dto.CreateServiceCategoryRequest,
	) (*dto.ServiceCategoryResponse, error)

	Update(
		ctx context.Context,
		id string,
		req dto.UpdateServiceCategoryRequest,
	) (*dto.ServiceCategoryResponse, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}
