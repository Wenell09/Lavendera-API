package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Wenell09/lavendera-api/internal/models"
	outletUserDto "github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/Wenell09/lavendera-api/internal/report/dto"
	"github.com/Wenell09/lavendera-api/internal/report/service"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type StubReportRepository struct {
	SummaryFunc            func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.DashboardSummary, error)
	OrderStatusRecapFunc   func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.OrderStatusRecap, error)
	PaymentMethodRecapFunc func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.PaymentMethodRecap, error)
	SalesChartFunc         func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.SalesChartItem, error)
	FinancialSummaryFunc   func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.FinancialSummary, error)
	PaymentsFunc           func(ctx context.Context, outletID *uuid.UUID, startDate, endDate time.Time, offset, limit int) ([]models.Payment, int64, error)
}

func (s *StubReportRepository) GetDashboardSummary(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.DashboardSummary, error) {
	if s.SummaryFunc != nil {
		return s.SummaryFunc(ctx, startDate, endDate, outletID)
	}
	return models.DashboardSummary{}, nil
}

func (s *StubReportRepository) GetOrderStatusRecap(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.OrderStatusRecap, error) {
	if s.OrderStatusRecapFunc != nil {
		return s.OrderStatusRecapFunc(ctx, startDate, endDate, outletID)
	}
	return nil, nil
}

func (s *StubReportRepository) GetPaymentMethodRecap(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.PaymentMethodRecap, error) {
	if s.PaymentMethodRecapFunc != nil {
		return s.PaymentMethodRecapFunc(ctx, startDate, endDate, outletID)
	}
	return nil, nil
}

func (s *StubReportRepository) GetSalesChart(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.SalesChartItem, error) {
	if s.SalesChartFunc != nil {
		return s.SalesChartFunc(ctx, startDate, endDate, outletID)
	}
	return nil, nil
}

func (s *StubReportRepository) GetFinancialSummary(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.FinancialSummary, error) {
	if s.FinancialSummaryFunc != nil {
		return s.FinancialSummaryFunc(ctx, startDate, endDate, outletID)
	}
	return models.FinancialSummary{}, nil
}

func (s *StubReportRepository) GetPayments(ctx context.Context, outletID *uuid.UUID, startDate, endDate time.Time, offset, limit int) ([]models.Payment, int64, error) {
	if s.PaymentsFunc != nil {
		return s.PaymentsFunc(ctx, outletID, startDate, endDate, offset, limit)
	}
	return nil, 0, nil
}

type StubOutletUserRepository struct {
	FindByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]models.OutletUser, error)
}

func (s *StubOutletUserRepository) FindByOutletID(ctx context.Context, outletID uuid.UUID, filter outletUserDto.OutletUserFilter) ([]models.OutletUser, int64, error) {
	return nil, 0, nil
}

func (s *StubOutletUserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.OutletUser, error) {
	if s.FindByUserIDFunc != nil {
		return s.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (s *StubOutletUserRepository) Exists(ctx context.Context, outletID, userID uuid.UUID) (bool, error) {
	return false, nil
}

func (s *StubOutletUserRepository) Create(ctx context.Context, outletUser *models.OutletUser) error {
	return nil
}

func (s *StubOutletUserRepository) Delete(ctx context.Context, outletID, userID uuid.UUID) error {
	return nil
}

func TestDashboard_Success(t *testing.T) {
	stubRepo := &StubReportRepository{
		SummaryFunc: func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.DashboardSummary, error) {
			return models.DashboardSummary{
				TotalSales:            1255000,
				TotalPaymentsReceived: 1322500,
				TotalOrders:           11,
				ActiveOrders:          3,
			}, nil
		},
		OrderStatusRecapFunc: func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.OrderStatusRecap, error) {
			return []models.OrderStatusRecap{
				{Status: "MENUNGGU_KONFIRMASI", Count: 1},
				{Status: "DIPROSES", Count: 2},
				{Status: "SELESAI", Count: 5},
				{Status: "DIAMBIL", Count: 2},
				{Status: "DIBATALKAN", Count: 1},
			}, nil
		},
	}
	stubOutletUser := &StubOutletUserRepository{}
	logger := logrus.New()

	svc := service.NewReportService(stubRepo, stubOutletUser, logger)

	tenantID := uuid.New()
	ctx := appcontext.WithTenantID(context.Background(), tenantID)
	ctx = appcontext.WithRole(ctx, "ADMIN")

	res, err := svc.GetDashboard(ctx, dto.DashboardFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Summary.TotalSales != 1255000 {
		t.Errorf("expected total_sales 1255000, got %d", res.Summary.TotalSales)
	}
	if res.Summary.TotalOrders != 11 {
		t.Errorf("expected total_orders 11, got %d", res.Summary.TotalOrders)
	}

	var sumStatus int64
	for _, s := range res.OrderStatus {
		sumStatus += s.Count
	}
	if res.Summary.TotalOrders != sumStatus {
		t.Errorf("total_orders %d does not match sum of status counts %d", res.Summary.TotalOrders, sumStatus)
	}
}

func TestFinancialReport_Reconciliation(t *testing.T) {
	outletID := uuid.New()
	orderID := uuid.New()
	payments := []models.Payment{
		{
			ID:            uuid.New(),
			Amount:        500000,
			Status:        "verified",
			OrderID:       orderID,
			OutletID:      outletID,
			PaymentMethod: "cash",
			Order: models.Order{
				OrderNumber: "ORD-001",
				Customer:    models.Customer{Name: "Customer A"},
			},
			Outlet: models.Outlet{Name: "Outlet A"},
		},
		{
			ID:            uuid.New(),
			Amount:        822500,
			Status:        "verified",
			OrderID:       orderID,
			OutletID:      outletID,
			PaymentMethod: "cash",
			Order: models.Order{
				OrderNumber: "ORD-001",
				Customer:    models.Customer{Name: "Customer A"},
			},
			Outlet: models.Outlet{Name: "Outlet A"},
		},
	}

	stubRepo := &StubReportRepository{
		FinancialSummaryFunc: func(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.FinancialSummary, error) {
			return models.FinancialSummary{
				TotalSales:            1255000,
				TotalDiscount:         10000,
				TotalPaymentsReceived: 1322500,
				TotalOrders:           11,
			}, nil
		},
		PaymentsFunc: func(ctx context.Context, outletID *uuid.UUID, startDate, endDate time.Time, offset, limit int) ([]models.Payment, int64, error) {
			return payments, int64(len(payments)), nil
		},
	}

	svc := service.NewReportService(stubRepo, &StubOutletUserRepository{}, logrus.New())

	tenantID := uuid.New()
	ctx := appcontext.WithTenantID(context.Background(), tenantID)
	ctx = appcontext.WithRole(ctx, "ADMIN")

	res, err := svc.GetFinancialReport(ctx, dto.FinancialReportFilter{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Data.Summary.TotalPaymentsReceived != 1322500 {
		t.Errorf("expected total_payments_received 1322500, got %d", res.Data.Summary.TotalPaymentsReceived)
	}

	var sumPayments int64
	for _, p := range res.Data.Payments {
		sumPayments += p.Amount
	}
	if res.Data.Summary.TotalPaymentsReceived != sumPayments {
		t.Errorf("total_payments_received %d does not match sum of payment list rows %d", res.Data.Summary.TotalPaymentsReceived, sumPayments)
	}
}
