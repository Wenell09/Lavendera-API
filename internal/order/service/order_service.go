package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/order/dto"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateCashierOrder(ctx context.Context, req dto.CreateCashierOrderRequest) (*dto.OrderResponse, error)
	CreatePublicOrder(ctx context.Context, req dto.CreatePublicOrderRequest) (*dto.OrderResponse, error)
	FindAll(ctx context.Context, filter dto.OrderFilter) (*dto.OrderListResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	FindByTrackingToken(ctx context.Context, trackingToken string) (*dto.OrderResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error)
	CancelOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	ConfirmPayment(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) (*dto.OrderResponse, error)
	InputCashPayment(ctx context.Context, id uuid.UUID, req dto.UpdateOrderPaidAmountRequest) (*dto.OrderResponse, error)
}
