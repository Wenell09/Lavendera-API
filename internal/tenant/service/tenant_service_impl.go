package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/Wenell09/lavendera-api/internal/tenant/dto"
	"github.com/Wenell09/lavendera-api/internal/tenant/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TenantServiceImpl struct {
	Repository repository.TenantRepository
	Logger     *logrus.Logger
}

func NewTenantService(repository repository.TenantRepository, logger *logrus.Logger) TenantService {
	return &TenantServiceImpl{Repository: repository, Logger: logger}
}

// FindByID implements [TenantService]
func (s *TenantServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.TenantResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("tenant_id", id)

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists || tenantID != id {
		return nil, apperror.UnauthorizedError{Msg: "unauthorized access to tenant"}
	}

	tenant, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "tenant not found"}
		}
		logger.WithError(err).Error("failed to find tenant")
		return nil, err
	}

	return &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Slug:      tenant.Slug,
		Email:     tenant.Email,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}, nil
}

// Update implements [TenantService]
func (s *TenantServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("tenant_id", id)

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists || tenantID != id {
		return nil, apperror.UnauthorizedError{Msg: "unauthorized access to tenant"}
	}

	tenant, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "tenant not found"}
		}
		logger.WithError(err).Error("failed to find tenant")
		return nil, err
	}

	updateData := make(map[string]interface{})
	if req.Name != nil {
		newName := strings.TrimSpace(*req.Name)
		if !strings.EqualFold(tenant.Name, newName) {
			exists, err := s.Repository.ExistsByName(ctx, newName)
			if err != nil {
				logger.WithField("new_name", newName).WithError(err).Error("failed to check tenant name existence")
				return nil, err
			}
			if exists {
				return nil, apperror.ConflictError{Msg: "tenant name already exists"}
			}
		}
		updateData["name"] = newName
	}
	if req.Slug != nil {
		newSlug := strings.TrimSpace(*req.Slug)
		if newSlug != tenant.Slug {
			exists, err := s.Repository.ExistsBySlug(ctx, newSlug)
			if err != nil {
				logger.WithField("new_slug", newSlug).WithError(err).Error("failed to check tenant slug existence")
				return nil, err
			}
			if exists {
				return nil, apperror.ConflictError{Msg: "tenant slug already exists"}
			}
		}
		updateData["slug"] = newSlug
	}

	if err := s.Repository.Update(ctx, id, updateData); err != nil {
		logger.WithError(err).Error("failed to update tenant")
		return nil, err
	}

	updated, _ := s.Repository.FindByID(ctx, id)
	return &dto.TenantResponse{
		ID:        updated.ID,
		Name:      updated.Name,
		Slug:      updated.Slug,
		Email:     updated.Email,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}
