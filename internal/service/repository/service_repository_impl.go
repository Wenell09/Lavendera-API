package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepositoryImpl struct {
	DB *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &ServiceRepositoryImpl{
		DB: db,
	}
}

// Create implements [ServiceRepository].
func (s *ServiceRepositoryImpl) Create(
	ctx context.Context,
	service *models.Service,
) error {
	return s.DB.WithContext(ctx).Create(service).Error
}

// FindAll implements [ServiceRepository].
func (s *ServiceRepositoryImpl) FindAll(
	ctx context.Context,
	search string,
	categoryID *uuid.UUID,
	page int,
	limit int,
) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64
	query := s.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Model(&models.Service{}).
		Preload("Category")

	if search != "" {
		query = query.Where(
			"name ILIKE ?",
			search+"%",
		)
	}
	if categoryID != nil {
		query = query.Where(
			"category_id = ?",
			*categoryID,
		)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&services).Error; err != nil {
		return nil, 0, err
	}
	return services, total, nil
}

// FindByID implements [ServiceRepository].
func (s *ServiceRepositoryImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Service, error) {
	var service models.Service
	err := s.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Preload("Category").
		First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &service, nil
}

// Update implements [ServiceRepository].
func (s *ServiceRepositoryImpl) Update(
	ctx context.Context,
	service *models.Service,
) error {
	return s.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Updates(service).Error
}

// Delete implements [ServiceRepository].
func (s *ServiceRepositoryImpl) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Where("id = ?", id).
		Delete(&models.Service{}).Error
}

// ExistByName implements [ServiceRepository].
func (s *ServiceRepositoryImpl) ExistByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := s.DB.
		WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Model(&models.Service{}).
		Where("name = ?", name).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
