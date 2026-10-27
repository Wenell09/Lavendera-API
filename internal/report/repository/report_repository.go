package repository

import (
	"context"
	"time"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/google/uuid"
)

type ReportRepository interface {
	GetDashboardSummary(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.DashboardSummary, error)
	GetOrderStatusRecap(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.OrderStatusRecap, error)
	GetPaymentMethodRecap(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.PaymentMethodRecap, error)
	GetSalesChart(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.SalesChartItem, error)

	GetFinancialSummary(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.FinancialSummary, error)
	GetPayments(ctx context.Context, outletID *uuid.UUID, startDate, endDate time.Time, offset, limit int) ([]models.Payment, int64, error)
}
