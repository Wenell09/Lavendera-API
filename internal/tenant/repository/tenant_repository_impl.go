package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepositoryImpl struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &TenantRepositoryImpl{DB: db}
}

// FindByID implements [TenantRepository].
func (r *TenantRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	tenant := &models.Tenant{}
	if err := r.DB.WithContext(ctx).First(tenant, id).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// Update implements [TenantRepository].
func (r *TenantRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	if err := r.DB.WithContext(ctx).Model(&models.Tenant{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByName implements [TenantRepository].
func (r *TenantRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&models.Tenant{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsBySlug implements [TenantRepository].
func (r *TenantRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&models.Tenant{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
