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

// Create implements [ServiceCategoryService].
func (s *ServiceCategoryServiceImpl) Create(ctx context.Context, req dto.CreateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error) {
	name := strings.TrimSpace(req.Name)
	s.Logger.WithField(
		"name",
		name,
	).Info("create service category started")
	exists, err := s.Repository.ExistsByName(
		ctx,
		name,
	)
	if err != nil {
		s.Logger.WithError(err).
			Error("failed to check service category name")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{
			Msg: "service category name already exists",
		}
	}
	category := &models.ServiceCategory{
		Name: name,
	}
	// tenant_id akan diisi dari JWT context.
	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{
			Msg: "tenant_id not found",
		}
	}
	category.TenantID, err = uuid.Parse(tenantID)
	if err != nil {
		return nil, apperror.UnauthorizedError{
			Msg: "invalid tenant_id",
		}
	}
	if err := s.Repository.Create(
		ctx,
		category,
	); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{
				Msg: "service category name already exists",
			}
		}
		s.Logger.WithError(err).
			Error("failed to create service category")
		return nil, err
	}
	s.Logger.WithField(
		"category_id",
		category.ID,
	).Info("service category created")
	return &dto.ServiceCategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Delete implements [ServiceCategoryService].
func (s *ServiceCategoryServiceImpl) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return apperror.ValidationError{
			Msg: "invalid service category id",
		}
	}
	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{
				Msg: "service category not found",
			}
		}
		return err
	}
	if err := s.Repository.Delete(
		ctx,
		category,
	); err != nil {
		s.Logger.WithError(err).
			Error("failed to delete service category")

		return err
	}
	s.Logger.WithField(
		"category_id",
		id,
	).Info("service category deleted")
	return nil
}

// FindAll implements [ServiceCategoryService].
func (s *ServiceCategoryServiceImpl) FindAll(ctx context.Context) (*dto.ServiceCategoryListResponse, error) {
	s.Logger.Info("find all service categories started")
	categories, err := s.Repository.FindAll(ctx)
	if err != nil {
		s.Logger.WithError(err).
			Error("failed to find service categories")
		return nil, err
	}
	data := make(
		[]dto.ServiceCategoryResponse,
		0,
		len(categories),
	)
	for _, category := range categories {
		data = append(
			data,
			dto.ServiceCategoryResponse{
				ID:        category.ID.String(),
				Name:      category.Name,
				CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		)
	}
	return &dto.ServiceCategoryListResponse{
		Data: data,
	}, nil
}

// FindByID implements [ServiceCategoryService].
func (s *ServiceCategoryServiceImpl) FindByID(ctx context.Context, id string) (*dto.ServiceCategoryResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperror.ValidationError{
			Msg: "invalid service category id",
		}
	}
	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{
				Msg: "service category not found",
			}
		}
		s.Logger.WithError(err).
			Error("failed to find service category")
		return nil, err
	}
	return &dto.ServiceCategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Update implements [ServiceCategoryService].
func (s *ServiceCategoryServiceImpl) Update(ctx context.Context, id string, req dto.UpdateServiceCategoryRequest) (*dto.ServiceCategoryResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperror.ValidationError{
			Msg: "invalid service category id",
		}
	}
	name := strings.TrimSpace(req.Name)
	category, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{
				Msg: "service category not found",
			}
		}
		return nil, err
	}
	// Jika nama tidak berubah, tidak perlu cek duplicate.
	if !strings.EqualFold(category.Name, name) {
		exists, err := s.Repository.ExistsByName(
			ctx,
			name,
		)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperror.ConflictError{
				Msg: "service category name already exists",
			}
		}
	}
	category.Name = name
	if err := s.Repository.Update(
		ctx,
		category,
	); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{
				Msg: "service category name already exists",
			}
		}
		s.Logger.WithError(err).
			Error("failed to update service category")

		return nil, err
	}
	s.Logger.WithField(
		"category_id",
		category.ID,
	).Info("service category updated")
	return &dto.ServiceCategoryResponse{
		ID:        category.ID.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
