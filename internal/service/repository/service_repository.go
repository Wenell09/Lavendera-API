package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/google/uuid"
)

type ServiceRepository interface {
	FindAll(ctx context.Context, filter dto.ServiceFilter) ([]models.Service, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Service, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, service *models.Service) error
	Update(ctx context.Context, id uuid.UUID, service map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
