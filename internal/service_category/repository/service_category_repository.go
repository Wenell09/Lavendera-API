package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type ServiceCategoryRepository interface {
	FindAll(
		ctx context.Context,
	) ([]models.ServiceCategory, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*models.ServiceCategory, error)

	ExistsByName(
		ctx context.Context,
		name string,
	) (bool, error)

	Create(
		ctx context.Context,
		category *models.ServiceCategory,
	) error

	Update(
		ctx context.Context,
		category *models.ServiceCategory,
	) error

	Delete(
		ctx context.Context,
		category *models.ServiceCategory,
	) error
}
