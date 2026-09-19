package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/google/uuid"
)

type Service interface {
	Create(
		ctx context.Context,
		request dto.CreateServiceRequest,
	) (*dto.ServiceResponse, error)

	FindAll(
		ctx context.Context,
		search string,
		categoryID *uuid.UUID,
		page int,
		limit int,
	) (*dto.ServiceListResponse, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*dto.ServiceResponse, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		request dto.UpdateServiceRequest,
	) (*dto.ServiceResponse, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
