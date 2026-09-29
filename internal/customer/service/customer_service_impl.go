package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/google/uuid"

	"github.com/Wenell09/lavendera-api/internal/customer/dto"
	"github.com/Wenell09/lavendera-api/internal/customer/repository"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerServiceImpl struct {
	Repository repository.CustomerRepository
	Logger     *logrus.Logger
}

func NewCustomerService(repository repository.CustomerRepository, logger *logrus.Logger) CustomerService {
	return &CustomerServiceImpl{Repository: repository, Logger: logger}
}

// Create implements [CustomerService].
func (c *CustomerServiceImpl) Create(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	logger := utils.LogWithContext(c.Logger, ctx).WithField("name", req.Name)
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	customer := &models.Customer{
		TenantID: tenantID,
		Name:     strings.TrimSpace(req.Name),
		Phone:    strings.TrimSpace(req.Phone),
		Address:  req.Address,
	}

	existsName, err := c.Repository.ExistsByName(ctx, customer.Name)
	if err != nil {
		logger.WithError(err).Error("failed to check customer name existence")
		return nil, err
	}
	if existsName {
		return nil, apperror.ConflictError{Msg: "customer name already exists"}
	}

	if err := c.Repository.Create(ctx, customer); err != nil {
		logger.WithError(err).Error("failed to create customer in database")
		return nil, err
	}

	logger.WithField("customer_id", customer.ID).Info("customer created successfully")
	return &dto.CustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		Phone:     customer.Phone,
		Address:   customer.Address,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}, nil
}

// Delete implements [CustomerService].
func (c *CustomerServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(c.Logger, ctx).WithField("customer_id", id)
	_, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "customer not found"}
		}
		logger.WithError(err).Error("failed to find customer for deletion")
		return err
	}

	if err := c.Repository.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete customer from database")
		return err
	}

	logger.Info("delete customer successfully")
	return nil
}

// FindAll implements [CustomerService].
func (c *CustomerServiceImpl) FindAll(ctx context.Context, filter dto.CustomerFilter) (*dto.CustomerListResponse, error) {
	filter.SetDefault()
	customers, total, err := c.Repository.FindAll(ctx, filter)
	if err != nil {
		utils.LogWithContext(c.Logger, ctx).WithFields(logrus.Fields{
			"search": filter.Search,
			"page":   filter.Page,
			"limit":  filter.Limit,
		}).WithError(err).Error("failed to fetch customers list")
		return nil, err
	}

	data := []dto.CustomerResponse{}
	for _, cust := range customers {
		data = append(data, dto.CustomerResponse{
			ID:        cust.ID,
			Name:      cust.Name,
			Phone:     cust.Phone,
			Address:   cust.Address,
			CreatedAt: cust.CreatedAt,
			UpdatedAt: cust.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	return &dto.CustomerListResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// FindByID implements [CustomerService].
func (c *CustomerServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error) {
	logger := utils.LogWithContext(c.Logger, ctx).WithField("customer_id", id)
	customer, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "customer not found"}
		}
		logger.WithError(err).Error("failed to find customer by id")
		return nil, err
	}

	return &dto.CustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		Phone:     customer.Phone,
		Address:   customer.Address,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}, nil
}

// Update implements [CustomerService].
func (c *CustomerServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	logger := utils.LogWithContext(c.Logger, ctx).WithField("customer_id", id)
	customer, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "customer not found"}
		}
		logger.WithError(err).Error("failed to find customer by id")
		return nil, err
	}

	updateData := make(map[string]interface{})
	if req.Name != nil {
		newName := strings.TrimSpace(*req.Name)
		if newName != customer.Name {
			existsName, err := c.Repository.ExistsByName(ctx, newName)
			if err != nil {
				logger.WithError(err).Error("failed to check customer name existence")
				return nil, err
			}
			if existsName {
				return nil, apperror.ConflictError{Msg: "customer name already exists"}
			}
		}
		updateData["name"] = newName
	}
	if req.Phone != nil {
		updateData["phone"] = strings.TrimSpace(*req.Phone)
	}
	if req.Address != nil {
		updateData["address"] = req.Address
	}

	if len(updateData) == 0 {
		return &dto.CustomerResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			Phone:     customer.Phone,
			Address:   customer.Address,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		}, nil
	}

	if err := c.Repository.Update(ctx, id, updateData); err != nil {
		logger.WithError(err).Error("failed to update customer in database")
		return nil, err
	}

	updatedCustomer, err := c.Repository.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch updated customer")
		return nil, err
	}

	logger.Info("customer updated successfully")
	return &dto.CustomerResponse{
		ID:        updatedCustomer.ID,
		Name:      updatedCustomer.Name,
		Phone:     updatedCustomer.Phone,
		Address:   updatedCustomer.Address,
		CreatedAt: updatedCustomer.CreatedAt,
		UpdatedAt: updatedCustomer.UpdatedAt,
	}, nil
}
