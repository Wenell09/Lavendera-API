package service

import (
	"context"
	"errors"
	"math"

	"github.com/Wenell09/lavendera-api/internal/models"
	outletRepo "github.com/Wenell09/lavendera-api/internal/outlet/repository"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	userRepo "github.com/Wenell09/lavendera-api/internal/user/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OutletUserServiceImpl struct {
	Repo       repository.OutletUserRepository
	OutletRepo outletRepo.OutletRepository
	UserRepo   userRepo.UserRepository
	Logger     *logrus.Logger
}

func NewOutletUserService(
	repo repository.OutletUserRepository,
	outletRepo outletRepo.OutletRepository,
	userRepo userRepo.UserRepository,
	logger *logrus.Logger,
) OutletUserService {
	return &OutletUserServiceImpl{
		Repo:       repo,
		OutletRepo: outletRepo,
		UserRepo:   userRepo,
		Logger:     logger,
	}
}

func (s *OutletUserServiceImpl) Assign(ctx context.Context, req dto.CreateOutletUserRequest) (*dto.OutletUserResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"outlet_id": req.OutletID,
		"user_id":   req.UserID,
	})

	if _, exists := appcontext.TenantIDFromContext(ctx); !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	outlet, err := s.OutletRepo.FindByID(ctx, req.OutletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet")
		return nil, err
	}

	user, err := s.UserRepo.FindByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "user not found"}
		}
		logger.WithError(err).Error("failed to find user")
		return nil, err
	}

	exists, err := s.Repo.Exists(ctx, req.OutletID, req.UserID)
	if err != nil {
		logger.WithError(err).Error("failed to check outlet_user existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "user already assigned to this outlet"}
	}

	outletUser := &models.OutletUser{
		OutletID: req.OutletID,
		UserID:   req.UserID,
	}

	if err := s.Repo.Create(ctx, outletUser); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "user already assigned to this outlet"}
		}
		logger.WithError(err).Error("failed to assign user to outlet")
		return nil, err
	}

	logger.Info("user assigned to outlet successfully")

	return &dto.OutletUserResponse{
		Outlet: dto.OutletUserOutletResponse{
			ID:   outlet.ID,
			Name: outlet.Name,
			Slug: outlet.Slug,
		},
		User: dto.OutletUserUserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
		CreatedAt: outletUser.CreatedAt,
	}, nil
}

func (s *OutletUserServiceImpl) Unassign(ctx context.Context, outletID, userID uuid.UUID) error {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"outlet_id": outletID,
		"user_id":   userID,
	})

	if _, exists := appcontext.TenantIDFromContext(ctx); !exists {
		return apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	if _, err := s.OutletRepo.FindByID(ctx, outletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet")
		return err
	}

	exists, err := s.Repo.Exists(ctx, outletID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check outlet_user existence")
		return err
	}
	if !exists {
		return apperror.NotFoundError{Msg: "user not assigned to this outlet"}
	}

	if err := s.Repo.Delete(ctx, outletID, userID); err != nil {
		logger.WithError(err).Error("failed to unassign user from outlet")
		return err
	}

	logger.Info("user unassigned from outlet successfully")
	return nil
}

func (s *OutletUserServiceImpl) FindByOutletID(ctx context.Context, outletID uuid.UUID, filter dto.OutletUserFilter) (*dto.OutletUserListResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("outlet_id", outletID)

	if _, exists := appcontext.TenantIDFromContext(ctx); !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	filter.SetDefault()

	outlet, err := s.OutletRepo.FindByID(ctx, outletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet")
		return nil, err
	}

	outletUsers, total, err := s.Repo.FindByOutletID(ctx, outletID, filter)
	if err != nil {
		logger.WithError(err).Error("failed to find outlet users")
		return nil, err
	}

	responses := []dto.OutletUserResponse{}
	for _, ou := range outletUsers {
		responses = append(responses, dto.OutletUserResponse{
			Outlet: dto.OutletUserOutletResponse{
				ID:   outlet.ID,
				Name: outlet.Name,
				Slug: outlet.Slug,
			},
			User: dto.OutletUserUserResponse{
				ID:       ou.User.ID,
				Name:     ou.User.Name,
				Email:    ou.User.Email,
				Role:     ou.User.Role,
				IsActive: ou.User.IsActive,
			},
			CreatedAt: ou.CreatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.OutletUserListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
