package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/Wenell09/lavendera-api/internal/user/dto"
	"github.com/Wenell09/lavendera-api/internal/user/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	Repository repository.UserRepository
	Logger     *logrus.Logger
}

func NewUserService(repository repository.UserRepository, logger *logrus.Logger) UserService {
	return &UserServiceImpl{Repository: repository, Logger: logger}
}

// Create implements [UserService].
func (u *UserServiceImpl) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)
	logger := utils.LogWithContext(u.Logger, ctx).WithFields(logrus.Fields{"email": email, "name": name})

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	exists, err := u.Repository.ExistsByEmail(ctx, email)
	if err != nil {
		logger.WithError(err).Error("failed to check user email existence")
		return nil, err
	}
	if exists {
		return nil, apperror.ConflictError{Msg: "email already registered"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("failed to hash password")
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := &models.User{
		TenantID: tenantID,
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "STAFF",
		IsActive: isActive,
	}

	if err := u.Repository.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "email already registered"}
		}
		logger.WithError(err).Error("failed to create user in database")
		return nil, err
	}

	logger.WithField("user_id", user.ID).Info("user created successfully")
	return &dto.UserResponse{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Delete implements [UserService].
func (u *UserServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(u.Logger, ctx).WithField("user_id", id)

	currentUserID, ok := appcontext.UserIDFromContext(ctx)
	if ok && currentUserID == id {
		return apperror.ConflictError{Msg: "cannot delete your own account"}
	}

	_, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "user not found"}
		}
		logger.WithError(err).Error("failed to find user for deletion")
		return err
	}

	if err := u.Repository.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete user from database")
		return err
	}

	logger.Info("user deleted successfully")
	return nil
}

// FindAll implements [UserService].
func (u *UserServiceImpl) FindAll(ctx context.Context, filter dto.UserFilter) (*dto.UserListResponse, error) {
	filter.SetDefault()
	users, total, err := u.Repository.FindAll(ctx, filter)
	if err != nil {
		utils.LogWithContext(u.Logger, ctx).WithFields(logrus.Fields{
			"search": filter.Search,
			"role":   filter.Role,
			"page":   filter.Page,
			"limit":  filter.Limit,
		}).WithError(err).Error("failed to fetch users list")
		return nil, err
	}

	data := []dto.UserResponse{}
	for _, user := range users {
		data = append(data, dto.UserResponse{
			ID:        user.ID,
			TenantID:  user.TenantID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	return &dto.UserListResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// FindByID implements [UserService].
func (u *UserServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	logger := utils.LogWithContext(u.Logger, ctx).WithField("user_id", id)
	user, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "user not found"}
		}
		logger.WithError(err).Error("failed to find user by id")
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Update implements [UserService].
func (u *UserServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	logger := utils.LogWithContext(u.Logger, ctx).WithField("user_id", id)

	user, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "user not found"}
		}
		logger.WithError(err).Error("failed to find user by id")
		return nil, err
	}

	updateData := make(map[string]interface{})
	if req.Name != nil {
		updateData["name"] = strings.TrimSpace(*req.Name)
	}

	if req.Email != nil {
		newEmail := strings.ToLower(strings.TrimSpace(*req.Email))
		if !strings.EqualFold(user.Email, newEmail) {
			exists, err := u.Repository.ExistsByEmail(ctx, newEmail)
			if err != nil {
				logger.WithField("new_email", newEmail).WithError(err).Error("failed to check email existence")
				return nil, err
			}
			if exists {
				return nil, apperror.ConflictError{Msg: "email already registered"}
			}
		}
		updateData["email"] = newEmail
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.WithError(err).Error("failed to hash password")
			return nil, err
		}
		updateData["password"] = string(hashedPassword)
	}

	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}

	if len(updateData) == 0 {
		return &dto.UserResponse{
			ID:        user.ID,
			TenantID:  user.TenantID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}, nil
	}

	if err := u.Repository.Update(ctx, id, updateData); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.ConflictError{Msg: "email already registered"}
		}
		logger.WithError(err).Error("failed to update user in database")
		return nil, err
	}

	updatedUser, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch updated user")
		return nil, err
	}

	logger.Info("user updated successfully")
	return &dto.UserResponse{
		ID:        updatedUser.ID,
		TenantID:  updatedUser.TenantID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
		IsActive:  updatedUser.IsActive,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}
