package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/google/uuid"
)

type OutletUserRepository interface {
	Create(ctx context.Context, outletUser *models.OutletUser) error
	Delete(ctx context.Context, outletID, userID uuid.UUID) error
	FindByOutletID(ctx context.Context, outletID uuid.UUID, filter dto.OutletUserFilter) ([]models.OutletUser, int64, error)
	Exists(ctx context.Context, outletID, userID uuid.UUID) (bool, error)
}
