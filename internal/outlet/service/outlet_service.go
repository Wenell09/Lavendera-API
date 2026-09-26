package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/outlet/dto"
	"github.com/google/uuid"
)

type OutletService interface {
	FindAll(ctx context.Context, filter dto.OutletFilter) (*dto.OutletListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.OutletResponse, error)
	Create(ctx context.Context, req dto.CreateOutletRequest) (*dto.OutletResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateOutletRequest) (*dto.OutletResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
