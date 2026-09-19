package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceCategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewServiceCategoryRepository(
	db *gorm.DB,
) ServiceCategoryRepository {
	return &ServiceCategoryRepositoryImpl{
		DB: db,
	}
}

// Create implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) Create(
	ctx context.Context,
	category *models.ServiceCategory,
) error {
	return s.DB.
		WithContext(ctx).
		Create(category).
		Error
}

// Delete implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) Delete(
	ctx context.Context,
	category *models.ServiceCategory,
) error {
	return s.DB.
		WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Where("id = ?", category.ID).
		Delete(&models.ServiceCategory{}).
		Error
}

// ExistsByName implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) ExistsByName(
	ctx context.Context,
	name string,
) (bool, error) {
	var count int64
	err := s.DB.
		WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Model(&models.ServiceCategory{}).
		Where("name = ?", name).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) FindAll(
	ctx context.Context,
) ([]models.ServiceCategory, error) {
	var categories []models.ServiceCategory
	err := s.DB.
		WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Order("created_at DESC").
		Find(&categories).
		Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.ServiceCategory, error) {
	category := &models.ServiceCategory{}
	err := s.DB.
		WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Where("id = ?", id).
		First(category).
		Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

// Update implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) Update(
	ctx context.Context,
	category *models.ServiceCategory,
) error {
	return s.DB.
		WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Where("id = ?", category.ID).
		Updates(category).
		Error
}
