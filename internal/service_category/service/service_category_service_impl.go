package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/service_category/dto"
	"github.com/Wenell09/lavendera-api/internal/service_category/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceCategoryServiceImpl struct {
	Repository repository.ServiceCategoryRepository
	Logger     *logrus.Logger
}

func NewServiceCategoryService(
	repository repository.ServiceCategoryRepository,
	logger *logrus.Logger,
) ServiceCategoryService {
	return &ServiceCategoryServiceImpl{
		Repository: repository,
		Logger:     logger,
	}
}

// Helper internal untuk menyertakan context default pada logger
func (s *ServiceCategoryServiceImpl) logWithCtx(ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(s.Logger)
	if tenantID, ok := appcontext.TenantIDFromContext(ctx); ok {
		entry = entry.WithField("tenant_id", tenantID)
	}
	return entry
}

func (s *ServiceCategoryServiceImpl) Create(ctx context.Context, req dto.CreateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error) {
	name := strings.TrimSpace(req.Name)
	logger := s.logWithCtx(ctx).WithField("name", name)
	exists, err := s.Repository.ExistsByName(ctx, name)
	if err != nil {
		logger.WithError(err).Error("failed to check service category name existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "service category name already exists"}
	}
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, apperror.UnauthorizedError{Msg: "invalid tenant_id"}
	}
	category := &models.ServiceCategory{
		Name:     name,
		TenantID: tenantUUID,
	}

	if err := s.Repository.Create(ctx, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "service category name already exists"}
		}
		logger.WithError(err).Error("failed to create service category in database")
		return nil, err
	}
	logger.WithField("category_id", category.ID).Info("service category created successfully")
	return &dto.ServiceCategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *ServiceCategoryServiceImpl) Delete(ctx context.Context, id string) error {
	logger := s.logWithCtx(ctx).WithField("category_id", id)
	if _, err := uuid.Parse(id); err != nil {
		return apperror.ValidationError{Msg: "invalid service category id"}
	}

	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "service category not found"}
		}
		logger.WithError(err).Error("failed to find service category for deletion")
		return err
	}
	if err := s.Repository.Delete(ctx, category); err != nil {
		logger.WithError(err).Error("failed to delete service category from database")
		return err
	}
	logger.Info("service category deleted successfully")
	return nil
}

func (s *ServiceCategoryServiceImpl) FindAll(ctx context.Context) (*dto.ServiceCategoryListResponse, error) {
	categories, err := s.Repository.FindAll(ctx)
	if err != nil {
		s.logWithCtx(ctx).WithError(err).Error("failed to find all service categories")
		return nil, err
	}
	data := make([]dto.ServiceCategoryResponse, 0, len(categories))
	for _, category := range categories {
		data = append(data, dto.ServiceCategoryResponse{
			ID:        category.ID.String(),
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &dto.ServiceCategoryListResponse{Data: data}, nil
}

func (s *ServiceCategoryServiceImpl) FindByID(ctx context.Context, id string) (*dto.ServiceCategoryResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperror.ValidationError{Msg: "invalid service category id"}
	}
	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "service category not found"}
		}
		s.logWithCtx(ctx).WithField("category_id", id).WithError(err).Error("failed to find service category by id")
		return nil, err
	}
	return &dto.ServiceCategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *ServiceCategoryServiceImpl) Update(ctx context.Context, id string, req dto.UpdateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error) {
	logger := s.logWithCtx(ctx).WithField("category_id", id)
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperror.ValidationError{Msg: "invalid service category id"}
	}
	name := strings.TrimSpace(req.Name)
	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "service category not found"}
		}
		logger.WithError(err).Error("failed to find service category for update")
		return nil, err
	}
	if !strings.EqualFold(category.Name, name) {
		exists, err := s.Repository.ExistsByName(ctx, name)
		if err != nil {
			logger.WithField("new_name", name).WithError(err).Error("failed to check service category name existence")
			return nil, err
		}
		if exists {
			return nil, apperror.ConflictError{Msg: "service category name already exists"}
		}
	}
	category.Name = name
	if err := s.Repository.Update(ctx, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "service category name already exists"}
		}
		logger.WithField("new_name", name).WithError(err).Error("failed to update service category in database")
		return nil, err
	}
	logger.Info("service category updated successfully")
	return &dto.ServiceCategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
