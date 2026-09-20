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
	return &AuthRepositoryImpl{DB: db}
}

// CreateTenantAndAdmin implements [AuthRepository].
func (a *AuthRepositoryImpl) CreateDefaultAdmin(ctx context.Context, tenant *models.Tenant, user *models.User, categories []models.ServiceCategory) error {
	return a.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}
		user.TenantID = tenant.ID
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		for i := range categories {
			categories[i].TenantID = tenant.ID
		}
		if err := tx.Create(categories).Error; err != nil {
			return err
		}
		return nil
	})
}

// ExistsTenantBySlug implements [AuthRepository].
func (a *AuthRepositoryImpl) ExistsTenantBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := a.DB.WithContext(ctx).Model(&models.Tenant{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsUserByEmail implements [AuthRepository].
func (a *AuthRepositoryImpl) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := a.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindUserByEmail implements [AuthRepository].
func (a *AuthRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	if err := a.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
