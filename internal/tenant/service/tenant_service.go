package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/tenant/dto"
	"github.com/google/uuid"
)

type TenantService interface {
	FindByID(ctx context.Context, id uuid.UUID) (*dto.TenantResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTenantRequest) (*dto.TenantResponse, error)
}
