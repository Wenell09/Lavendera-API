package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	FindAll(ctx context.Context, search string, categoryID *uuid.UUID, page int, limit int) ([]models.Service, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Service, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistByName(ctx context.Context, name string) (bool, error)
}
