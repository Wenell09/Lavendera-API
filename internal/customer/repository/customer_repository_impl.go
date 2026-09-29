package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/customer/dto"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	DB *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{DB: db}
}

// Create implements [CustomerRepository].
func (c *CustomerRepositoryImpl) Create(ctx context.Context, customer *models.Customer) error {
	if err := c.DB.WithContext(ctx).Create(customer).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements [CustomerRepository].
func (c *CustomerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Delete(&models.Customer{}, id).Error; err != nil {
		return err
	}
	return nil
}

// FindAll implements [CustomerRepository].
func (c *CustomerRepositoryImpl) FindAll(ctx context.Context, filter dto.CustomerFilter) ([]models.Customer, int64, error) {
	customers := []models.Customer{}
	var total int64
	query := c.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Customer{})
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR phone ILIKE ?", filter.Search+"%", filter.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&customers).Error; err != nil {
		return nil, 0, err
	}
	return customers, total, nil
}

// FindByID implements [CustomerRepository].
func (c *CustomerRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	customer := &models.Customer{}
	if err := c.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).First(customer, id).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

// ExistsByName implements [CustomerRepository].
func (c *CustomerRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := c.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Customer{}).Where("name ILIKE ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update implements [CustomerRepository].
func (c *CustomerRepositoryImpl) Update(ctx context.Context, id uuid.UUID, customer map[string]interface{}) error {
	if len(customer) == 0 {
		return nil
	}
	if err := c.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Customer{}).Where("id = ?", id).Updates(customer).Error; err != nil {
		return err
	}
	return nil
}
