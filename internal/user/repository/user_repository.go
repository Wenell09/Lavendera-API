package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/user/dto"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindAll(ctx context.Context, filter dto.UserFilter) ([]models.User, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id uuid.UUID, user map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
