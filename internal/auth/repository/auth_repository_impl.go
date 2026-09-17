package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

// CreateTenantAndAdmin implements [AuthRepository].
func (a *AuthRepositoryImpl) CreateDefaultAdmin(ctx context.Context, tenant *models.Tenant, user *models.User, categories []*models.ServiceCategory) error {
	return a.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create tenant
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}
		// Set tenant ID untuk admin
		user.TenantID = tenant.ID
		// Create admin
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		for _, value := range categories {
			value.TenantID = tenant.ID
		}
		// Create category default
		if err := tx.Create(categories).Error; err != nil {
			return err
		}
		return nil
	})
}

// ExistsTenantBySlug implements [AuthRepository].
func (a *AuthRepositoryImpl) ExistsTenantBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := a.DB.WithContext(ctx).
		Model(&models.Tenant{}).
		Where("slug = ?", slug).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsUserByEmail implements [AuthRepository].
func (a *AuthRepositoryImpl) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := a.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindUserByEmail implements [AuthRepository].
func (a *AuthRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := a.DB.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
