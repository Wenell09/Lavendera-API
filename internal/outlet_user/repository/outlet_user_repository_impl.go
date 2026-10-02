package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutletUserRepositoryImpl struct {
	DB *gorm.DB
}

func NewOutletUserRepository(db *gorm.DB) OutletUserRepository {
	return &OutletUserRepositoryImpl{DB: db}
}

// Create implements [OutletUserRepository].
func (o *OutletUserRepositoryImpl) Create(ctx context.Context, outletUser *models.OutletUser) error {
	if err := o.DB.WithContext(ctx).Create(outletUser).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [OutletUserRepository].
func (o *OutletUserRepositoryImpl) Delete(ctx context.Context, outletID, userID uuid.UUID) error {
	if err := o.DB.WithContext(ctx).Scopes(database.TenantScopeByOutlet(ctx)).
		Where("outlet_id = ? AND user_id = ?", outletID, userID).
		Delete(&models.OutletUser{}).Error; err != nil {
		return err
	}
	return nil
}

// Exists implements [OutletUserRepository].
func (o *OutletUserRepositoryImpl) Exists(ctx context.Context, outletID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := o.DB.WithContext(ctx).Scopes(database.TenantScopeByOutlet(ctx)).
		Model(&models.OutletUser{}).
		Where("outlet_id = ? AND user_id = ?", outletID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByOutletID implements [OutletUserRepository].
func (o *OutletUserRepositoryImpl) FindByOutletID(ctx context.Context, outletID uuid.UUID, filter dto.OutletUserFilter) ([]models.OutletUser, int64, error) {
	outletUsers := []models.OutletUser{}
	var total int64
	query := o.DB.WithContext(ctx).Scopes(database.TenantScopeByOutlet(ctx)).
		Model(&models.OutletUser{}).
		Where("outlet_users.outlet_id = ?", outletID)
	if filter.Search != "" {
		searchPattern := filter.Search + "%"
		query = query.Joins("JOIN users ON users.id = outlet_users.user_id").
			Where("users.name ILIKE ? OR users.email ILIKE ?", searchPattern, searchPattern)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Preload("User").Preload("Outlet").Order("outlet_users.created_at DESC").Offset(offset).Limit(filter.Limit).Find(&outletUsers).Error; err != nil {
		return nil, 0, err
	}
	return outletUsers, total, nil
}
