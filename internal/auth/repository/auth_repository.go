package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
)

type AuthRepository interface {
	CreateDefaultAdmin(
		ctx context.Context,
		tenant *models.Tenant,
		user *models.User,
		category []*models.ServiceCategory,
	) error

	FindUserByEmail(
		ctx context.Context,
		email string,
	) (*models.User, error)

	ExistsUserByEmail(
		ctx context.Context,
		email string,
	) (bool, error)

	ExistsTenantBySlug(
		ctx context.Context,
		slug string,
	) (bool, error)
}
