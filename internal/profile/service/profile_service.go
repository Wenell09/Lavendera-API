package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/profile/dto"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (*dto.ProfileResponse, error)
}
