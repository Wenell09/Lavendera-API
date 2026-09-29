package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/google/uuid"
)

type OutletUserService interface {
	Assign(ctx context.Context, req dto.CreateOutletUserRequest) (*dto.OutletUserResponse, error)
	Unassign(ctx context.Context, outletID, userID uuid.UUID) error
	FindByOutletID(ctx context.Context, outletID uuid.UUID, filter dto.OutletUserFilter) (*dto.OutletUserListResponse, error)
}
