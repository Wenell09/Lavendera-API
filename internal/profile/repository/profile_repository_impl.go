package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepositoryImpl struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &ProfileRepositoryImpl{
		DB: db,
	}
}

// GetProfile implements [ProfileRepository].
func (p *ProfileRepositoryImpl) GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	profile := &models.User{}
	if err := p.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).First(profile, id).Error; err != nil {
		return nil, err
	}
	return profile, nil
}
