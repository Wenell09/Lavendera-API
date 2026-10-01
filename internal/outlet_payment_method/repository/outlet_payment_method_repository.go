package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/dto"
	"github.com/google/uuid"
)

type OutletPaymentMethodRepository interface {
	FindAll(ctx context.Context, outletID uuid.UUID, filter dto.OutletPaymentMethodFilter) ([]models.OutletPaymentMethod, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.OutletPaymentMethod, error)
	Exists(ctx context.Context, outletID uuid.UUID, paymentType string, providerName string, accountNumber *string) (bool, error)
	Create(ctx context.Context, outletPaymentMethod *models.OutletPaymentMethod) error
	Update(ctx context.Context, id uuid.UUID, paymentMethod map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
