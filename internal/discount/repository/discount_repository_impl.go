package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/discount/dto"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountRepositoryImpl struct {
	DB *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &DiscountRepositoryImpl{DB: db}
}

// Create implements [DiscountRepository].
func (d *DiscountRepositoryImpl) Create(ctx context.Context, discount *models.Discount) error {
	if err := d.DB.WithContext(ctx).Create(discount).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [DiscountRepository].
func (d *DiscountRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := d.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Delete(&models.Discount{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByName implements [DiscountRepository].
func (d *DiscountRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := d.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Discount{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll implements [DiscountRepository].
func (d *DiscountRepositoryImpl) FindAll(ctx context.Context, filter dto.DiscountFilter) ([]models.Discount, int64, error) {
	discounts := []models.Discount{}
	var total int64
	query := d.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Discount{})
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", filter.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&discounts).Error; err != nil {
		return nil, 0, err
	}
	return discounts, total, nil
}

// FindByID implements [DiscountRepository].
func (d *DiscountRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Discount, error) {
	discount := &models.Discount{}
	if err := d.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).First(discount, id).Error; err != nil {
		return nil, err
	}
	return discount, nil
}

// Update implements [DiscountRepository].
func (d *DiscountRepositoryImpl) Update(ctx context.Context, id uuid.UUID, discount map[string]interface{}) error {
	if len(discount) == 0 {
		return nil
	}
	if err := d.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Discount{}).Where("id = ?", id).Updates(discount).Error; err != nil {
		return err
	}
	return nil
}
