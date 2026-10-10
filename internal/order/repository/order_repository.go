package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/order/dto"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateWithPayment(ctx context.Context, order *models.Order, payment *models.Payment) error
	FindAll(ctx context.Context, filter dto.OrderFilter) ([]models.Order, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	FindByTrackingToken(ctx context.Context, trackingToken string) (*models.Order, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	ExistsByTrackingToken(ctx context.Context, trackingToken string) (bool, error)
	ExistsByOrderNumber(ctx context.Context, tenantID uuid.UUID, orderNumber string) (bool, error)
	FindPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	SumVerifiedPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) (int64, error)
	UpdatePaymentAndOrderPaidAmount(ctx context.Context, paymentID uuid.UUID, paymentUpdates map[string]interface{}, orderID uuid.UUID, newPaidAmount int64) error
	CancelOrderWithPayments(ctx context.Context, orderID uuid.UUID, orderUpdates map[string]interface{}, paymentUpdates map[string]interface{}) error
}
