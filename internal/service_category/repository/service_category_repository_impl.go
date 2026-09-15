package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/database"
	"github.com/Wenell09/lavendera-api/internal/models"
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
func (s *ServiceCategoryRepositoryImpl) Create(ctx context.Context, category *models.ServiceCategory) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return gorm.ErrInvalidDB
	}
	return tx.
		WithContext(ctx).
		Create(category).
		Error
}

// Delete implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) Delete(ctx context.Context, category *models.ServiceCategory) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return gorm.ErrInvalidDB
	}
	return tx.
		WithContext(ctx).
		Delete(category).
		Error
}

// ExistsByName implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return false, gorm.ErrInvalidDB
	}
	var count int64
	if err := tx.
		WithContext(ctx).
		Model(&models.ServiceCategory{}).
		Where("name = ?", name).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) FindAll(ctx context.Context) ([]models.ServiceCategory, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return nil, gorm.ErrInvalidDB
	}
	var categories []models.ServiceCategory
	if err := tx.
		WithContext(ctx).
		Order("created_at DESC").
		Find(&categories).
		Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) FindByID(ctx context.Context, id string) (*models.ServiceCategory, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return nil, gorm.ErrInvalidDB
	}
	category := &models.ServiceCategory{}
	if err := tx.
		WithContext(ctx).
		Where("id = ?", id).
		First(category).
		Error; err != nil {
		return nil, err
	}
	return category, nil
}

// Update implements [ServiceCategoryRepository].
func (s *ServiceCategoryRepositoryImpl) Update(ctx context.Context, category *models.ServiceCategory) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return gorm.ErrInvalidDB
	}
	return tx.
		WithContext(ctx).
		Save(category).
		Error
}
