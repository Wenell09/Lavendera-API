package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/Wenell09/lavendera-api/internal/discount/dto"
)

type DiscountService interface {
	FindAll(ctx context.Context, filter dto.DiscountFilter) (*dto.DiscountListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.DiscountResponse, error)
	Create(ctx context.Context, req dto.CreateDiscountRequest) (*dto.DiscountResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateDiscountRequest) (*dto.DiscountResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
