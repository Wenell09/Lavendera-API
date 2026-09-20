package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepositoryImpl struct {
	DB *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &ServiceRepositoryImpl{DB: db}
}

// Create implements [ServiceRepository].
func (s *ServiceRepositoryImpl) Create(ctx context.Context, service *models.Service) error {
	if err := s.DB.WithContext(ctx).Create(service).Error; err != nil {
		return err
	}
	return nil
}

// FindAll implements [ServiceRepository].
func (s *ServiceRepositoryImpl) FindAll(ctx context.Context, filter dto.ServiceFilter) ([]models.Service, int64, error) {
	services := []models.Service{}
	var total int64
	query := s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Service{}).Preload("Category")
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", filter.Search+"%")
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&services).Error; err != nil {
		return nil, 0, err
	}
	return services, total, nil
}

// FindByID implements [ServiceRepository].
func (s *ServiceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	service := &models.Service{}
	if err := s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Preload("Category").First(&service, id).Error; err != nil {
		return nil, err
	}
	return service, nil
}

// Update implements [ServiceRepository].
func (s *ServiceRepositoryImpl) Update(ctx context.Context, id uuid.UUID, service map[string]interface{}) error {
	if len(service) == 0 {
		return nil
	}
	if err := s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Service{}).Where("id = ?", id).Updates(service).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [ServiceRepository].
func (s *ServiceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Delete(&models.Service{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ExistByName implements [ServiceRepository].
func (s *ServiceRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Service{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
