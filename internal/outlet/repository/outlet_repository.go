package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet/dto"
	"github.com/google/uuid"
)

type OutletRepository interface {
	FindAll(ctx context.Context, filter dto.OutletFilter) ([]models.Outlet, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Outlet, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, outlet *models.Outlet) error
	Update(ctx context.Context, id uuid.UUID, outlet map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
