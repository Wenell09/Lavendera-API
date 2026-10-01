package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/dto"
	"github.com/google/uuid"
)

type OutletPaymentMethodService interface {
	Create(ctx context.Context, req dto.CreateOutletPaymentMethodRequest) (*dto.OutletPaymentMethodResponse, error)
	FindAll(ctx context.Context, outletID uuid.UUID, filter dto.OutletPaymentMethodFilter) (*dto.OutletPaymentMethodListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.OutletPaymentMethodResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateOutletPaymentMethodRequest) (*dto.OutletPaymentMethodResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
