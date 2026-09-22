package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/Wenell09/lavendera-api/internal/service/repository"
	categoryRepository "github.com/Wenell09/lavendera-api/internal/service_category/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	ServiceRepository  repository.ServiceRepository
	CategoryRepository categoryRepository.ServiceCategoryRepository
	Logger             *logrus.Logger
}

func NewService(serviceRepository repository.ServiceRepository, logger *logrus.Logger, categoryRepository categoryRepository.ServiceCategoryRepository) Service {
	return &ServiceImpl{ServiceRepository: serviceRepository, Logger: logger, CategoryRepository: categoryRepository}
}

// Create implements [Service].
func (s *ServiceImpl) Create(ctx context.Context, req dto.CreateServiceRequest) (*dto.ServiceResponse, error) {
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}
	exists, err := s.ServiceRepository.ExistsByName(ctx, req.Name)
	if err != nil {
		utils.LogWithContext(s.Logger, ctx).WithError(err).Error("failed to check service category name existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "service category name already exists"}
	}
	service := &models.Service{
		TenantID:     tenantID,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Price:        req.Price,
		MinQuantity:  req.MinQuantity,
		Unit:         req.Unit,
		DurationDays: req.DurationDays,
		IsActive:     req.IsActive,
	}
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"service_name": service.Name,
		"category_id":  service.CategoryID,
	})
	if err := s.ServiceRepository.Create(ctx, service); err != nil {
		logger.WithError(err).Error("failed to create service in database")
		return nil, err
	}
	createdService, err := s.ServiceRepository.FindByID(ctx, service.ID)
	if err != nil {
		logger.WithField("service_id", service.ID).WithError(err).Error("failed to fetch created service details")
		return nil, err
	}
	logger.WithField("service_id", service.ID).Info("service created successfully")
	response := mapServiceResponse(createdService)
	return &response, nil
}

// FindAll implements [Service].
func (s *ServiceImpl) FindAll(ctx context.Context, filter dto.ServiceFilter) (*dto.ServiceListResponse, error) {
	filter.SetDefault()
	services, total, err := s.ServiceRepository.FindAll(ctx, filter)
	if err != nil {
		utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
			"search":      filter.Search,
			"category_id": filter.CategoryID,
			"page":        filter.Page,
			"limit":       filter.Limit,
		}).WithError(err).Error("failed to fetch services list")
		return nil, err
	}
	data := []dto.ServiceResponse{}
	for i := range services {
		data = append(data, mapServiceResponse(&services[i]))
	}
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	return &dto.ServiceListResponse{
		Data:       data,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// FindByID implements [Service].
func (s *ServiceImpl) FindByID(ctx context.Context, id uuid.UUID,
) (*dto.ServiceResponse, error) {
	service, err := s.ServiceRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "service not found"}
		}
		utils.LogWithContext(s.Logger, ctx).WithField("service_id", id).WithError(err).Error("failed to find service by id")
		return nil, err
	}
	response := mapServiceResponse(service)
	return &response, nil
}

// Update implements [Service].
func (s *ServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateServiceRequest) (*dto.ServiceResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("service_id", id)
	service, err := s.ServiceRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "service not found"}
		}
		logger.WithError(err).Error("failed to find service for update")
		return nil, err
	}
	if req.CategoryID != nil {
		_, err := s.CategoryRepository.FindByID(ctx, *req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperror.NotFoundError{Msg: "service category not found"}
			}
			logger.WithError(err).Error("failed to find service category")
			return nil, err
		}
	}
	updateData := make(map[string]interface{})
	if req.Name != nil {
		newName := strings.TrimSpace(*req.Name)
		if !strings.EqualFold(service.Name, newName) {
			exists, err := s.ServiceRepository.ExistsByName(ctx, newName)
			if err != nil {
				logger.WithField("new_name", newName).WithError(err).Error("failed to check service name existence")
				return nil, err
			}
			if exists {
				return nil, apperror.ConflictError{Msg: "service name already exists"}
			}
		}
		updateData["name"] = newName
	}
	if req.CategoryID != nil {
		updateData["category_id"] = *req.CategoryID
	}
	if req.Price != nil {
		updateData["price"] = *req.Price
	}
	if req.MinQuantity != nil {
		updateData["min_quantity"] = *req.MinQuantity
	}
	if req.Unit != nil {
		updateData["unit"] = *req.Unit
	}
	if req.DurationDays != nil {
		updateData["duration_days"] = *req.DurationDays
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}
	if len(updateData) == 0 {
		response := mapServiceResponse(service)
		return &response, nil
	}
	if err := s.ServiceRepository.Update(ctx, id, updateData); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "service name already exists"}
		}
		logger.WithError(err).Error("failed to update service in database")
		return nil, err
	}
	updatedService, err := s.ServiceRepository.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch updated service")
		return nil, err
	}
	logger.Info("service updated successfully")
	response := mapServiceResponse(updatedService)
	return &response, nil
}

// Delete implements [Service].
func (s *ServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("service_id", id)
	_, err := s.ServiceRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "service not found"}
		}
		logger.WithError(err).Error("failed to find service for deletion")
		return err
	}
	if err := s.ServiceRepository.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete service from database")
		return err
	}
	logger.Info("service deleted successfully")
	return nil
}

func mapServiceResponse(service *models.Service) dto.ServiceResponse {
	return dto.ServiceResponse{
		ID: service.ID,
		Category: dto.CategoryResponse{
			ID:   service.CategoryID,
			Name: service.Category.Name,
		},
		Name:         service.Name,
		Price:        service.Price,
		MinQuantity:  service.MinQuantity,
		Unit:         service.Unit,
		DurationDays: service.DurationDays,
		IsActive:     service.IsActive,
		CreatedAt:    service.CreatedAt,
		UpdatedAt:    service.UpdatedAt,
	}
}
