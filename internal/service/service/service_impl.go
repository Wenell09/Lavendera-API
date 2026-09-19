package service

import (
	"context"
	"errors"
	"math"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/Wenell09/lavendera-api/internal/service/repository"
	categoryRepository "github.com/Wenell09/lavendera-api/internal/service_category/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	ServiceRepository  repository.ServiceRepository
	CategoryRepository categoryRepository.ServiceCategoryRepository
	Logger             *logrus.Logger
}

func NewService(
	serviceRepository repository.ServiceRepository,
	logger *logrus.Logger,
	categoryRepository categoryRepository.ServiceCategoryRepository,
) Service {
	return &ServiceImpl{
		ServiceRepository:  serviceRepository,
		Logger:             logger,
		CategoryRepository: categoryRepository,
	}
}

// Helper internal untuk menginjeksikan context ID (seperti tenant_id) ke log entry
func (s *ServiceImpl) logCtx(ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(s.Logger)
	if tenantID, ok := appcontext.TenantIDFromContext(ctx); ok {
		entry = entry.WithField("tenant_id", tenantID)
	}
	return entry
}

// Create implements [Service].
func (s *ServiceImpl) Create(
	ctx context.Context,
	request dto.CreateServiceRequest,
) (*dto.ServiceResponse, error) {
	exists, err := s.ServiceRepository.ExistByName(ctx, request.Name)
	if err != nil {
		s.logCtx(ctx).WithError(err).Error("failed to check service category name existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "service category name already exists"}
	}
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}
	isActive := true
	if request.IsActive != nil {
		isActive = *request.IsActive
	}
	service := &models.Service{
		TenantID:     tenantID,
		CategoryID:   request.CategoryID,
		Name:         request.Name,
		Price:        request.Price,
		MinQuantity:  request.MinQuantity,
		Unit:         request.Unit,
		DurationDays: request.DurationDays,
		IsActive:     isActive,
	}

	logger := s.logCtx(ctx).WithFields(logrus.Fields{
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
func (s *ServiceImpl) FindAll(
	ctx context.Context,
	search string,
	categoryID *uuid.UUID,
	page int,
	limit int,
) (*dto.ServiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	services, total, err := s.ServiceRepository.FindAll(
		ctx,
		search,
		categoryID,
		page,
		limit,
	)
	if err != nil {
		s.logCtx(ctx).WithFields(logrus.Fields{
			"search":      search,
			"category_id": categoryID,
			"page":        page,
			"limit":       limit,
		}).WithError(err).Error("failed to fetch services list")
		return nil, err
	}

	data := []dto.ServiceResponse{}
	for i := range services {
		data = append(data, mapServiceResponse(&services[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.ServiceListResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// FindByID implements [Service].
func (s *ServiceImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.ServiceResponse, error) {
	service, err := s.ServiceRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "service not found"}
		}
		s.logCtx(ctx).WithField("service_id", id).WithError(err).Error("failed to find service by id")
		return nil, err
	}
	response := mapServiceResponse(service)
	return &response, nil
}

// Update implements [Service].
func (s *ServiceImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	request dto.UpdateServiceRequest,
) (*dto.ServiceResponse, error) {
	logger := s.logCtx(ctx).WithField("service_id", id)
	service, err := s.ServiceRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{
				Msg: "service not found",
			}
		}
		logger.WithError(err).
			Error("failed to find service for update")
		return nil, err
	}
	if request.CategoryID != nil {
		category, err := s.CategoryRepository.FindByID(
			ctx,
			*request.CategoryID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperror.NotFoundError{
					Msg: "service category not found",
				}
			}
			logger.WithError(err).
				Error("failed to find service category")
			return nil, err
		}
		service.CategoryID = *request.CategoryID
		service.Category = *category
	}
	updateServiceFields(service, request)
	if err := s.ServiceRepository.Update(ctx, service); err != nil {
		logger.WithError(err).
			Error("failed to update service in database")
		return nil, err
	}
	logger.Info("service updated successfully")
	response := mapServiceResponse(service)
	return &response, nil
}

// Delete implements [Service].
func (s *ServiceImpl) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	logger := s.logCtx(ctx).WithField("service_id", id)
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
		CreatedAt:    service.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    service.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func updateServiceFields(
	service *models.Service,
	request dto.UpdateServiceRequest,
) {
	if request.Name != nil {
		service.Name = *request.Name
	}
	if request.Price != nil {
		service.Price = *request.Price
	}
	if request.MinQuantity != nil {
		service.MinQuantity = *request.MinQuantity
	}
	if request.Unit != nil {
		service.Unit = *request.Unit
	}
	if request.DurationDays != nil {
		service.DurationDays = *request.DurationDays
	}
	if request.IsActive != nil {
		service.IsActive = *request.IsActive
	}
}
