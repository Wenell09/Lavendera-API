package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type TenantRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	Update(ctx context.Context, id uuid.UUID, data map[string]interface{}) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}
