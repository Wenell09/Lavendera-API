package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/customer/dto"
	"github.com/google/uuid"
)

type CustomerService interface {
	FindAll(ctx context.Context, filter dto.CustomerFilter) (*dto.CustomerListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error)
	Create(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
