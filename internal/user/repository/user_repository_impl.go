package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/Wenell09/lavendera-api/internal/user/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

// Create implements [UserRepository].
func (u *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [UserRepository].
func (u *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByEmail implements [UserRepository].
func (u *UserRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := u.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll implements [UserRepository].
func (u *UserRepositoryImpl) FindAll(ctx context.Context, filter dto.UserFilter) ([]models.User, int64, error) {
	users := []models.User{}
	var total int64
	query := u.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.User{})
	if filter.Search != "" {
		searchPattern := filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// FindByID implements [UserRepository].
func (u *UserRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := u.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Update implements [UserRepository].
func (u *UserRepositoryImpl) Update(ctx context.Context, id uuid.UUID, user map[string]interface{}) error {
	if len(user) == 0 {
		return nil
	}
	if err := u.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}
	return nil
}
