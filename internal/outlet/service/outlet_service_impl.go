package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet/dto"
	"github.com/Wenell09/lavendera-api/internal/outlet/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OutletServiceImpl struct {
	Repository repository.OutletRepository
	Logger     *logrus.Logger
}

func NewOutletService(repository repository.OutletRepository, logger *logrus.Logger) OutletService {
	return &OutletServiceImpl{Repository: repository, Logger: logger}
}

// Create implements [OutletService].
func (o *OutletServiceImpl) Create(ctx context.Context, req dto.CreateOutletRequest) (*dto.OutletResponse, error) {
	logger := utils.LogWithContext(o.Logger, ctx).WithField("name", req.Name)
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}
	exists, err := o.Repository.ExistsByName(ctx, req.Name)
	if err != nil {
		logger.WithError(err).Error("failed to check outlet name existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "outlet name already exists"}
	}
	outlet := &models.Outlet{
		TenantID:             tenantID,
		Name:                 req.Name,
		Slug:                 req.Slug,
		Phone:                req.Phone,
		Address:              req.Address,
		IsPublicOrderEnabled: req.IsPublicOrderEnabled,
		IsActive:             req.IsActive,
	}
	if err := o.Repository.Create(ctx, outlet); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "outlet name already exists"}
		}
		logger.WithError(err).Error("failed to create outlet in database")
		return nil, err
	}
	logger.WithField("outlet_id", outlet.ID).Info("outlet create successfully")
	return &dto.OutletResponse{
		ID:                   outlet.ID,
		Name:                 outlet.Name,
		Slug:                 outlet.Slug,
		Phone:                outlet.Phone,
		Address:              outlet.Address,
		IsPublicOrderEnabled: outlet.IsPublicOrderEnabled,
		IsActive:             outlet.IsActive,
		CreatedAt:            outlet.CreatedAt,
		UpdatedAt:            outlet.UpdatedAt,
	}, nil
}

// Delete implements [OutletService].
func (o *OutletServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(o.Logger, ctx).WithField("outlet_id", id)
	_, err := o.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet for deletion")
		return err
	}
	if err := o.Repository.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete outlet from database")
		return err
	}
	logger.Info("delete outlet successfully")
	return nil
}

// FindAll implements [OutletService].
func (o *OutletServiceImpl) FindAll(ctx context.Context, filter dto.OutletFilter) (*dto.OutletListResponse, error) {
	filter.SetDefault()
	outlets, total, err := o.Repository.FindAll(ctx, filter)
	if err != nil {
		utils.LogWithContext(o.Logger, ctx).WithFields(logrus.Fields{
			"search": filter.Search,
			"page":   filter.Page,
			"limit":  filter.Limit,
		}).WithError(err).Error("failed to fetch outlets list")
		return nil, err
	}
	data := []dto.OutletResponse{}
	for _, outlet := range outlets {
		data = append(data, dto.OutletResponse{
			ID:                   outlet.ID,
			Name:                 outlet.Name,
			Slug:                 outlet.Slug,
			Phone:                outlet.Phone,
			Address:              outlet.Address,
			IsPublicOrderEnabled: outlet.IsPublicOrderEnabled,
			IsActive:             outlet.IsActive,
			CreatedAt:            outlet.CreatedAt,
			UpdatedAt:            outlet.UpdatedAt,
		})
	}
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	return &dto.OutletListResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// FindByID implements [OutletService].
func (o *OutletServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.OutletResponse, error) {
	logger := utils.LogWithContext(o.Logger, ctx).WithField("outlet_id", id)
	outlet, err := o.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet by id")
		return nil, err
	}
	return &dto.OutletResponse{
		ID:                   outlet.ID,
		Name:                 outlet.Name,
		Slug:                 outlet.Slug,
		Phone:                outlet.Phone,
		Address:              outlet.Address,
		IsPublicOrderEnabled: outlet.IsPublicOrderEnabled,
		IsActive:             outlet.IsActive,
		CreatedAt:            outlet.CreatedAt,
		UpdatedAt:            outlet.UpdatedAt,
	}, nil
}

// Update implements [OutletService].
func (o *OutletServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateOutletRequest) (*dto.OutletResponse, error) {
	logger := utils.LogWithContext(o.Logger, ctx).WithField("outlet_id", id)
	outlet, err := o.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet by id")
		return nil, err
	}
	updateData := make(map[string]interface{})
	if req.Name != nil {
		newName := strings.TrimSpace(*req.Name)
		if !strings.EqualFold(outlet.Name, newName) {
			exists, err := o.Repository.ExistsByName(ctx, newName)
			if err != nil {
				logger.WithField("new_name", newName).WithError(err).Error("failed to check outlet name existence")
				return nil, err
			}
			if exists {
				return nil, apperror.ConflictError{Msg: "outlet name already exists"}
			}
		}
		updateData["name"] = newName
	}
	if req.Slug != nil {
		updateData["slug"] = *req.Slug
	}
	if req.Phone != nil {
		updateData["phone"] = *req.Phone
	}
	if req.Address != nil {
		updateData["address"] = *req.Address
	}
	if req.IsPublicOrderEnabled != nil {
		updateData["is_public_order_enabled"] = *req.IsPublicOrderEnabled
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}
	if len(updateData) == 0 {
		return &dto.OutletResponse{
			ID:                   outlet.ID,
			Name:                 outlet.Name,
			Slug:                 outlet.Slug,
			Phone:                outlet.Phone,
			Address:              outlet.Address,
			IsPublicOrderEnabled: outlet.IsPublicOrderEnabled,
			IsActive:             outlet.IsActive,
			CreatedAt:            outlet.CreatedAt,
			UpdatedAt:            outlet.UpdatedAt,
		}, nil
	}
	if err := o.Repository.Update(ctx, id, updateData); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "outlet name already exists"}
		}
		logger.WithError(err).Error("failed to update outlet in database")
		return nil, err
	}
	updatedOutlet, err := o.Repository.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch updated outlet")
		return nil, err
	}
	logger.Info("outlet updated successfully")
	return &dto.OutletResponse{
		ID:                   updatedOutlet.ID,
		Name:                 updatedOutlet.Name,
		Slug:                 updatedOutlet.Slug,
		Phone:                updatedOutlet.Phone,
		Address:              updatedOutlet.Address,
		IsPublicOrderEnabled: updatedOutlet.IsPublicOrderEnabled,
		IsActive:             updatedOutlet.IsActive,
		CreatedAt:            updatedOutlet.CreatedAt,
		UpdatedAt:            updatedOutlet.UpdatedAt,
	}, nil
}
