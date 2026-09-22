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
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceCategoryServiceImpl struct {
	Repository repository.ServiceCategoryRepository
	Logger     *logrus.Logger
}

func NewServiceCategoryService(repository repository.ServiceCategoryRepository, logger *logrus.Logger) ServiceCategoryService {
	return &ServiceCategoryServiceImpl{Repository: repository, Logger: logger}
}

func (s *ServiceCategoryServiceImpl) Create(ctx context.Context, req dto.CreateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error) {
	name := strings.TrimSpace(req.Name)
	logger := utils.LogWithContext(s.Logger, ctx).WithField("name", name)
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}
	exists, err := s.Repository.ExistsByName(ctx, name)
	if err != nil {
		logger.WithError(err).Error("failed to check service category name existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "service category name already exists"}
	}
	category := &models.ServiceCategory{
		Name:     name,
		TenantID: tenantID,
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
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (s *ServiceCategoryServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("category_id", id)
	_, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "service category not found"}
		}
		logger.WithError(err).Error("failed to find service category for deletion")
		return err
	}
	if err := s.Repository.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete service category from database")
		return err
	}
	logger.Info("service category deleted successfully")
	return nil
}

func (s *ServiceCategoryServiceImpl) FindAll(ctx context.Context) ([]dto.ServiceCategoryResponse, error) {
	categories, err := s.Repository.FindAll(ctx)
	if err != nil {
		utils.LogWithContext(s.Logger, ctx).WithError(err).Error("failed to find all service categories")
		return nil, err
	}
	data := []dto.ServiceCategoryResponse{}
	for _, category := range categories {
		data = append(data, dto.ServiceCategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
	}
	return data, nil
}

func (s *ServiceCategoryServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.ServiceCategoryResponse, error) {
	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "service category not found"}
		}
		utils.LogWithContext(s.Logger, ctx).WithField("category_id", id).WithError(err).Error("failed to find service category by id")
		return nil, err
	}
	return &dto.ServiceCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (s *ServiceCategoryServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("category_id", id)
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
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}
