package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/customer/dto"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	FindAll(ctx context.Context, filter dto.CustomerFilter) ([]models.Customer, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Customer, error)
	Create(ctx context.Context, customer *models.Customer) error
	Update(ctx context.Context, id uuid.UUID, customer map[string]interface{}) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
