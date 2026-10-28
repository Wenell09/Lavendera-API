package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error)
}
