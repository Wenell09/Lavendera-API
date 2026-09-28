package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
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

func (r *OutletUserRepositoryImpl) FindByOutletID(ctx context.Context, outletID uuid.UUID) ([]models.OutletUser, error) {
	var outletUsers []models.OutletUser
	err := r.DB.WithContext(ctx).
		Preload("User").
		Preload("Outlet").
		Where("outlet_id = ?", outletID).
		Order("created_at DESC").
		Find(&outletUsers).Error
	return outletUsers, err
}

func (r *OutletUserRepositoryImpl) Exists(ctx context.Context, outletID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.OutletUser{}).
		Where("outlet_id = ? AND user_id = ?", outletID, userID).
		Count(&count).Error
	return count > 0, err
}
