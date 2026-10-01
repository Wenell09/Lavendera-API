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

// Create implements [OutletPaymentMethodRepository].
func (o *OutletPaymentMethodRepositoryImpl) Create(ctx context.Context, outletPaymentMethod *models.OutletPaymentMethod) error {
	if err := o.DB.WithContext(ctx).Create(outletPaymentMethod).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [OutletPaymentMethodRepository].
func (o *OutletPaymentMethodRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := o.DB.WithContext(ctx).Scopes(database.OutletPaymentMethodsTenantScope(ctx)).
		Where("id = ?", id).Delete(&models.OutletPaymentMethod{}).Error; err != nil {
		return err
	}
	return nil
}

// Exists implements [OutletPaymentMethodRepository].
func (o *OutletPaymentMethodRepositoryImpl) Exists(ctx context.Context, outletID uuid.UUID, paymentType string, providerName string, accountNumber *string) (bool, error) {
	var count int64
	query := o.DB.WithContext(ctx).Scopes(database.OutletPaymentMethodsTenantScope(ctx)).Model(&models.OutletPaymentMethod{}).
		Where("outlet_id = ? AND type = ? AND provider_name = ?", outletID, paymentType, providerName)

	if accountNumber != nil {
		query = query.Where("account_number = ?", *accountNumber)
	} else {
		query = query.Where("account_number IS NULL")
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll implements [OutletPaymentMethodRepository].
func (o *OutletPaymentMethodRepositoryImpl) FindAll(ctx context.Context, outletID uuid.UUID, filter dto.OutletPaymentMethodFilter) ([]models.OutletPaymentMethod, int64, error) {
	payments := []models.OutletPaymentMethod{}
	var total int64

	query := o.DB.WithContext(ctx).Scopes(database.OutletPaymentMethodsTenantScope(ctx)).Model(&models.OutletPaymentMethod{}).
		Where("outlet_id = ?", outletID)

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("provider_name ILIKE ? OR account_name ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Preload("Outlet").Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}

// FindByID implements [OutletPaymentMethodRepository].
func (o *OutletPaymentMethodRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.OutletPaymentMethod, error) {
	payment := &models.OutletPaymentMethod{}
	if err := o.DB.WithContext(ctx).Scopes(database.OutletPaymentMethodsTenantScope(ctx)).Preload("Outlet").
		Where("id = ?", id).First(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

// Update implements [OutletPaymentMethodRepository].
func (o *OutletPaymentMethodRepositoryImpl) Update(ctx context.Context, id uuid.UUID, paymentMethod map[string]interface{}) error {
	if len(paymentMethod) == 0 {
		return nil
	}
	if err := o.DB.WithContext(ctx).Scopes(database.OutletPaymentMethodsTenantScope(ctx)).
		Model(&models.OutletPaymentMethod{}).Where("id = ?", id).Updates(paymentMethod).Error; err != nil {
		return err
	}
	return nil
}
