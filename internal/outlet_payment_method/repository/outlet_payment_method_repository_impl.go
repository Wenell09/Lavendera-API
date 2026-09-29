package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/dto"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutletPaymentMethodRepositoryImpl struct {
	DB *gorm.DB
}

func NewOutletPaymentMethodRepository(db *gorm.DB) OutletPaymentMethodRepository {
	return &OutletPaymentMethodRepositoryImpl{DB: db}
}

func (r *OutletPaymentMethodRepositoryImpl) Exists(ctx context.Context, outletID uuid.UUID, paymentType string, providerName string, accountNumber *string) (bool, error) {
	var count int64
	query := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.OutletPaymentMethod{}).
		Where("outlet_id = ? AND type = ? AND provider_name = ?", outletID, paymentType, providerName)

	if accountNumber != nil {
		query = query.Where("account_number = ?", *accountNumber)
	} else {
		query = query.Where("account_number IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *OutletPaymentMethodRepositoryImpl) Create(ctx context.Context, outletPaymentMethod *models.OutletPaymentMethod) error {
	return r.DB.WithContext(ctx).Create(outletPaymentMethod).Error
}

func (r *OutletPaymentMethodRepositoryImpl) FindAll(ctx context.Context, outletID uuid.UUID, filter dto.OutletPaymentMethodFilter) ([]models.OutletPaymentMethod, int64, error) {
	var payments []models.OutletPaymentMethod
	var total int64

	query := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.OutletPaymentMethod{}).Where("outlet_id = ?", outletID)

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("provider_name ILIKE ? OR account_name ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Preload("Outlet").Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&payments).Error
	return payments, total, err
}

func (r *OutletPaymentMethodRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.OutletPaymentMethod, error) {
	var payment models.OutletPaymentMethod
	err := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Preload("Outlet").Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *OutletPaymentMethodRepositoryImpl) Update(ctx context.Context, outletPaymentMethod *models.OutletPaymentMethod) error {
	return r.DB.WithContext(ctx).Save(outletPaymentMethod).Error
}

func (r *OutletPaymentMethodRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Where("id = ?", id).Delete(&models.OutletPaymentMethod{}).Error
}
