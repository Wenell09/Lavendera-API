package service

import (
	"context"
	"errors"

	"github.com/Wenell09/lavendera-api/internal/profile/dto"
	"github.com/Wenell09/lavendera-api/internal/profile/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProfileServiceImpl struct {
	Repository repository.ProfileRepository
	Logger     *logrus.Logger
}

func NewProfileService(repository repository.ProfileRepository, logger *logrus.Logger) ProfileService {
	return &ProfileServiceImpl{
		Repository: repository,
		Logger:     logger,
	}
}

// GetProfile implements [ProfileService].
func (p *ProfileServiceImpl) GetProfile(ctx context.Context) (*dto.ProfileResponse, error) {
	id, exist := appcontext.UserIDFromContext(ctx)
	logger := utils.LogWithContext(p.Logger, ctx).WithField("user_id", id)
	if !exist {
		return nil, apperror.UnauthorizedError{Msg: "User ID not found"}
	}
	profile, err := p.Repository.GetProfile(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "Profile not found"}
		}
		logger.WithError(err).Error("failed to find profile data")
		return nil, err
	}
	return &dto.ProfileResponse{
		ID:        profile.ID,
		Name:      profile.Name,
		Email:     profile.Email,
		Role:      profile.Role,
		CreatedAt: profile.CreatedAt,
	}, nil
}
