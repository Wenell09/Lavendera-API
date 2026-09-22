package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/google/uuid"

	"github.com/Wenell09/lavendera-api/internal/discount/dto"
	"github.com/Wenell09/lavendera-api/internal/discount/repository"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DiscountServiceImpl struct {
	Repository repository.DiscountRepository
	Logger     *logrus.Logger
}

func NewDiscountService(repository repository.DiscountRepository, logger *logrus.Logger) DiscountService {
	return &DiscountServiceImpl{Repository: repository, Logger: logger}
}

// Create implements [DiscountService].
func (d *DiscountServiceImpl) Create(ctx context.Context, req dto.CreateDiscountRequest) (*dto.DiscountResponse, error) {
	logger := utils.LogWithContext(d.Logger, ctx).WithField("name", req.Name)
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}
	exists, err := d.Repository.ExistsByName(ctx, req.Name)
	if err != nil {
		logger.WithError(err).Error("failed to check discount name existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "discount name already exists"}
	}
	discount := &models.Discount{
		TenantID: tenantID,
		Name:     req.Name,
		Type:     req.Type,
		Value:    req.Value,
		IsActive: req.IsActive,
	}
	if err := d.Repository.Create(ctx, discount); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "discount name already exists"}
		}
		logger.WithError(err).Error("failed to create discount in database")
		return nil, err
	}
	logger.WithField("discount_id", discount.ID).Info("discount create successfully")
	return &dto.DiscountResponse{
		ID:        discount.ID,
		Name:      discount.Name,
		Type:      discount.Type,
		Value:     discount.Value,
		IsActive:  discount.IsActive,
		CreatedAt: discount.CreatedAt,
		UpdatedAt: discount.UpdatedAt,
	}, nil
}

// Delete implements [DiscountService].
func (d *DiscountServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(d.Logger, ctx).WithField("discount_id", id)
	_, err := d.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "discount not found"}
		}
		logger.WithError(err).Error("failed to find discount for deletion")
		return err
	}
	if err := d.Repository.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete discount from database")
		return err
	}
	logger.Info("delete discount successfully")
	return nil
}

// FindAll implements [DiscountService].
func (d *DiscountServiceImpl) FindAll(ctx context.Context, filter dto.DiscountFilter) (*dto.DiscountListResponse, error) {
	filter.SetDefault()
	discounts, total, err := d.Repository.FindAll(ctx, filter)
	if err != nil {
		utils.LogWithContext(d.Logger, ctx).WithFields(logrus.Fields{
			"search": filter.Search,
			"page":   filter.Page,
			"limit":  filter.Limit,
		}).WithError(err).Error("failed to fetch discounts list")
		return nil, err
	}
	data := []dto.DiscountResponse{}
	for _, discount := range discounts {
		data = append(data, dto.DiscountResponse{
			ID:        discount.ID,
			Name:      discount.Name,
			Type:      discount.Type,
			Value:     discount.Value,
			IsActive:  discount.IsActive,
			CreatedAt: discount.CreatedAt,
			UpdatedAt: discount.UpdatedAt,
		})
	}
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	return &dto.DiscountListResponse{
		Data:       data,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// FindByID implements [DiscountService].
func (d *DiscountServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.DiscountResponse, error) {
	logger := utils.LogWithContext(d.Logger, ctx).WithField("discount_id", id)
	discount, err := d.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "discount not found"}
		}
		logger.WithError(err).Error("failed to find discount by id")
		return nil, err
	}
	return &dto.DiscountResponse{
		ID:        discount.ID,
		Name:      discount.Name,
		Type:      discount.Type,
		Value:     discount.Value,
		IsActive:  discount.IsActive,
		CreatedAt: discount.CreatedAt,
		UpdatedAt: discount.UpdatedAt,
	}, nil
}

// Update implements [DiscountService].
func (d *DiscountServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateDiscountRequest) (*dto.DiscountResponse, error) {
	logger := utils.LogWithContext(d.Logger, ctx).WithField("discount_id", id)
	discount, err := d.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "discount not found"}
		}
		logger.WithError(err).Error("failed to find discount by id")
		return nil, err
	}
	updateData := make(map[string]interface{})
	if req.Name != nil {
		newName := strings.TrimSpace(*req.Name)
		if !strings.EqualFold(discount.Name, newName) {
			exists, err := d.Repository.ExistsByName(ctx, newName)
			if err != nil {
				logger.WithField("new_name", newName).WithError(err).Error("failed to check discount name existence")
				return nil, err
			}
			if exists {
				return nil, apperror.ConflictError{Msg: "discount name already exists"}
			}
		}
		updateData["name"] = newName
	}
	if req.Type != nil {
		updateData["type"] = *req.Type
	}
	if req.Value != nil {
		updateData["value"] = *req.Value
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}
	if len(updateData) == 0 {
		return &dto.DiscountResponse{
			ID:        discount.ID,
			Name:      discount.Name,
			Type:      discount.Type,
			Value:     discount.Value,
			IsActive:  discount.IsActive,
			CreatedAt: discount.CreatedAt,
			UpdatedAt: discount.UpdatedAt,
		}, nil
	}
	if err := d.Repository.Update(ctx, id, updateData); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "discount name already exists"}
		}
		logger.WithError(err).Error("failed to update discount in database")
		return nil, err
	}
	updatedDiscount, err := d.Repository.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch updated discount")
		return nil, err
	}
	logger.Info("discount updated successfully")
	return &dto.DiscountResponse{
		ID:        updatedDiscount.ID,
		Name:      updatedDiscount.Name,
		Type:      updatedDiscount.Type,
		Value:     updatedDiscount.Value,
		IsActive:  updatedDiscount.IsActive,
		CreatedAt: updatedDiscount.CreatedAt,
		UpdatedAt: updatedDiscount.UpdatedAt,
	}, nil
}
