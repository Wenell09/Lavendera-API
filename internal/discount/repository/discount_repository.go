package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/discount/dto"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type DiscountRepository interface {
	FindAll(ctx context.Context, filter dto.DiscountFilter) ([]models.Discount, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Discount, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, discount *models.Discount) error
	Update(ctx context.Context, id uuid.UUID, discount map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
