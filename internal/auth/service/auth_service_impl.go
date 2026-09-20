package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/auth/dto"
	"github.com/Wenell09/lavendera-api/internal/auth/repository"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/config"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	Repository repository.AuthRepository
	JWTConfig  config.JWTConfig
	Logger     *logrus.Logger
}

func NewAuthService(repository repository.AuthRepository, jwtConfig config.JWTConfig, logger *logrus.Logger) AuthService {
	return &AuthServiceImpl{Repository: repository, JWTConfig: jwtConfig, Logger: logger}
}

func (a *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	a.Logger.WithField("email", req.Email).Info("login request started")
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := a.Repository.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.Logger.WithField("email", email).Warn("login failed: invalid credentials")
			return nil, apperror.UnauthorizedError{Msg: "invalid email or password"}
		}
		a.Logger.WithError(err).Error("failed to find user")
		return nil, err
	}
	if !user.IsActive {
		a.Logger.WithField("user_id", user.ID).Warn("login failed: user inactive")
		return nil, apperror.UnauthorizedError{Msg: "user account is inactive"}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		a.Logger.WithField("email", email).Warn("login failed: invalid credentials")
		return nil, apperror.UnauthorizedError{Msg: "invalid email or password"}
	}
	token, err := utils.GenerateToken(user.ID.String(), user.TenantID.String(), a.JWTConfig)
	if err != nil {
		a.Logger.WithError(err).Error("failed to generate JWT")
		return nil, err
	}
	a.Logger.WithFields(logrus.Fields{"user_id": user.ID, "tenant_id": user.TenantID}).Info("login successful")
	return &dto.LoginResponse{
		User: dto.UserResponse{
			ID:       user.ID.String(),
			TenantID: user.TenantID.String(),
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
		},
		Token: token,
	}, nil
}

func (a *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	a.Logger.WithFields(logrus.Fields{"email": req.Email}).Info("register request started")
	email := strings.ToLower(strings.TrimSpace(req.Email))
	slug := utils.GenerateSlug(req.TenantName)
	emailExists, err := a.Repository.ExistsUserByEmail(ctx, email)
	if err != nil {
		a.Logger.WithError(err).Error("failed to check email availability")
		return nil, err
	}
	if emailExists {
		a.Logger.WithField("email", email).Warn("register failed: email already registered")
		return nil, apperror.ConflictError{Msg: "email already registered"}
	}
	slugExists, err := a.Repository.ExistsTenantBySlug(ctx, slug)
	if err != nil {
		a.Logger.WithError(err).Error("failed to check tenant slug availability")
		return nil, err
	}
	if slugExists {
		a.Logger.WithField("slug", slug).Warn("register failed: tenant slug already exists")
		return nil, apperror.ConflictError{Msg: "tenant slug already exists"}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		a.Logger.WithError(err).Error("failed to hash password")
		return nil, err
	}
	tenant := &models.Tenant{Name: req.TenantName, Slug: slug, Email: email}
	user := &models.User{Name: req.Name, Email: email, Password: string(hashedPassword), Role: "ADMIN", IsActive: true}
	categories := []models.ServiceCategory{{Name: "Kiloan"}, {Name: "Satuan"}}
	if err := a.Repository.CreateDefaultAdmin(ctx, tenant, user, categories); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			a.Logger.WithError(err).Warn("register failed: duplicated data")
			return nil, apperror.ConflictError{Msg: "email or tenant slug already exists"}
		}
		a.Logger.WithError(err).Error("failed to create tenant and owner")
		return nil, err
	}
	a.Logger.WithFields(logrus.Fields{"tenant_id": tenant.ID, "user_id": user.ID, "email": email}).Info("register successful")
	return &dto.RegisterResponse{
		Tenant: dto.TenantResponse{ID: tenant.ID.String(), Name: tenant.Name, Slug: tenant.Slug, Email: tenant.Email},
		User:   dto.UserResponse{ID: user.ID.String(), TenantID: user.TenantID.String(), Name: user.Name, Email: user.Email, Role: user.Role},
	}, nil
}
