package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutletUserRepositoryImpl struct {
	DB *gorm.DB
}

func NewOutletUserRepository(db *gorm.DB) OutletUserRepository {
	return &OutletUserRepositoryImpl{DB: db}
}

func (r *OutletUserRepositoryImpl) Create(ctx context.Context, outletUser *models.OutletUser) error {
	return r.DB.WithContext(ctx).Create(outletUser).Error
}

func (r *OutletUserRepositoryImpl) Delete(ctx context.Context, outletID, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Where("outlet_id = ? AND user_id = ?", outletID, userID).
		Delete(&models.OutletUser{}).Error
}

func (r *OutletUserRepositoryImpl) FindByOutletID(ctx context.Context, outletID uuid.UUID, filter dto.OutletUserFilter) ([]models.OutletUser, int64, error) {
	var outletUsers []models.OutletUser
	var total int64
	query := r.DB.WithContext(ctx).
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
	err := query.Preload("User").Preload("Outlet").Order("outlet_users.created_at DESC").Offset(offset).Limit(filter.Limit).Find(&outletUsers).Error
	return outletUsers, total, err
}

func (r *OutletUserRepositoryImpl) Exists(ctx context.Context, outletID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.OutletUser{}).
		Where("outlet_id = ? AND user_id = ?", outletID, userID).
		Count(&count).Error
	return count > 0, err
}
