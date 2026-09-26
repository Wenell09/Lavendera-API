package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet/dto"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutletRepositoryImpl struct {
	DB *gorm.DB
}

func NewOutletRepository(db *gorm.DB) OutletRepository { return &OutletRepositoryImpl{DB: db} }

// Create implements [OutletRepository].
func (o *OutletRepositoryImpl) Create(ctx context.Context, outlet *models.Outlet) error {
	if err := o.DB.WithContext(ctx).Create(outlet).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [OutletRepository].
func (o *OutletRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := o.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Delete(&models.Outlet{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByName implements [OutletRepository].
func (o *OutletRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := o.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Outlet{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll implements [OutletRepository].
func (o *OutletRepositoryImpl) FindAll(ctx context.Context, filter dto.OutletFilter) ([]models.Outlet, int64, error) {
	outlets := []models.Outlet{}
	var total int64
	query := o.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Outlet{})
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", filter.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&outlets).Error; err != nil {
		return nil, 0, err
	}
	return outlets, total, nil
}

// FindByID implements [OutletRepository].
func (o *OutletRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Outlet, error) {
	outlet := &models.Outlet{}
	if err := o.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).First(outlet, id).Error; err != nil {
		return nil, err
	}
	return outlet, nil
}

// Update implements [OutletRepository].
func (o *OutletRepositoryImpl) Update(ctx context.Context, id uuid.UUID, outlet map[string]interface{}) error {
	if len(outlet) == 0 {
		return nil
	}
	if err := o.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Outlet{}).Where("id = ?", id).Updates(outlet).Error; err != nil {
		return err
	}
	return nil
}
