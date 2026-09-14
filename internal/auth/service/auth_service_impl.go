package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Wenell09/lavendera-api/internal/auth/dto"
	"github.com/Wenell09/lavendera-api/internal/auth/repository"
	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

type AuthServiceImpl struct {
	Repository repository.AuthRepository
	JWTConfig  config.JWTConfig
	Logger     *logrus.Logger
}

func NewAuthService(
	repository repository.AuthRepository,
	jwtConfig config.JWTConfig,
	logger *logrus.Logger,
) AuthService {
	return &AuthServiceImpl{
		Repository: repository,
		JWTConfig:  jwtConfig,
		Logger:     logger,
	}
}

// Login implements [AuthService].
func (a *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	a.Logger.WithField(
		"email",
		req.Email,
	).Info("login request started")
	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)
	// Find user
	user, err := a.Repository.FindUserByEmail(
		ctx,
		email,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.Logger.WithField(
				"email",
				email,
			).Warn("login failed: invalid credentials")
			return nil, utils.UnauthorizedError{
				Msg: "invalid email or password",
			}
		}
		a.Logger.WithError(err).
			Error("failed to find user")
		return nil, err
	}
	// Check active user
	if !user.IsActive {
		a.Logger.WithField(
			"user_id",
			user.ID,
		).Warn("login failed: user inactive")

		return nil, utils.UnauthorizedError{
			Msg: "user account is inactive",
		}
	}
	// Compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		a.Logger.WithField(
			"email",
			email,
		).Warn("login failed: invalid credentials")
		return nil, utils.UnauthorizedError{
			Msg: "invalid email or password",
		}
	}
	// Generate JWT
	token, err := GenerateToken(
		user.ID.String(),
		user.TenantID.String(),
		a.JWTConfig,
	)
	if err != nil {
		a.Logger.WithError(err).
			Error("failed to generate JWT")

		return nil, err
	}

	a.Logger.WithFields(logrus.Fields{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
	}).Info("login successful")

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

// Register implements [AuthService].
func (a *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	a.Logger.WithFields(logrus.Fields{
		"email": req.Email,
		"slug":  req.Slug,
	}).Info("register request started")

	email := strings.ToLower(strings.TrimSpace(req.Email))
	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	// Check email
	emailExists, err := a.Repository.ExistsUserByEmail(
		ctx,
		email,
	)
	if err != nil {
		a.Logger.WithError(err).
			Error("failed to check email availability")
		return nil, err
	}
	if emailExists {
		a.Logger.WithField("email", email).
			Warn("register failed: email already registered")
		return nil, utils.ConflictError{
			Msg: "email already registered",
		}
	}
	// Check tenant slug
	slugExists, err := a.Repository.ExistsTenantBySlug(
		ctx,
		slug,
	)
	if err != nil {
		a.Logger.WithError(err).
			Error("failed to check tenant slug availability")
		return nil, err
	}
	if slugExists {
		a.Logger.WithField("slug", slug).
			Warn("register failed: tenant slug already exists")
		return nil, utils.ConflictError{
			Msg: "tenant slug already exists",
		}
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		a.Logger.WithError(err).
			Error("failed to hash password")
		return nil, err
	}
	// Create models
	tenant := &models.Tenant{
		Name:  req.TenantName,
		Slug:  slug,
		Email: email,
	}
	user := &models.User{
		Name:     req.Name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "ADMIN",
		IsActive: true,
	}
	// Create tenant + admin
	if err := a.Repository.CreateTenantAndAdmin(
		ctx,
		tenant,
		user,
	); err != nil {
		// Handle race condition:
		// dua request register email yang sama
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			a.Logger.WithError(err).
				Warn("register failed: duplicated data")
			return nil, utils.ConflictError{
				Msg: "email or tenant slug already exists",
			}
		}
		a.Logger.WithError(err).
			Error("failed to create tenant and owner")
		return nil, err
	}
	a.Logger.WithFields(logrus.Fields{
		"tenant_id": tenant.ID,
		"user_id":   user.ID,
		"email":     email,
	}).Info("register successful")
	return &dto.RegisterResponse{
		Tenant: dto.TenantResponse{
			ID:    tenant.ID.String(),
			Name:  tenant.Name,
			Slug:  tenant.Slug,
			Email: tenant.Email,
		},
		User: dto.UserResponse{
			ID:       user.ID.String(),
			TenantID: user.TenantID.String(),
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}
