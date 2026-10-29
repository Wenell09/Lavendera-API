package service

import (
	"context"
	"math"
	"time"

	"github.com/Wenell09/lavendera-api/internal/models"
	outletUserRepo "github.com/Wenell09/lavendera-api/internal/outlet_user/repository"
	"github.com/Wenell09/lavendera-api/internal/report/dto"
	"github.com/Wenell09/lavendera-api/internal/report/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var defaultLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	defaultLocation = loc
}

type ReportServiceImpl struct {
	Repository     repository.ReportRepository
	OutletUserRepo outletUserRepo.OutletUserRepository
	Logger         *logrus.Logger
}

func NewReportService(
	repo repository.ReportRepository,
	outletUserRepo outletUserRepo.OutletUserRepository,
	logger *logrus.Logger,
) ReportService {
	return &ReportServiceImpl{
		Repository:     repo,
		OutletUserRepo: outletUserRepo,
		Logger:         logger,
	}
}

func (s *ReportServiceImpl) GetDashboard(ctx context.Context, filter dto.DashboardFilter) (*dto.DashboardResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx)

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists || tenantID == uuid.Nil {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	if filter.OutletID != nil {
		if err := s.validateOutletAccess(ctx, *filter.OutletID); err != nil {
			return nil, err
		}
	}

	startDate, endDate, err := resolveDateRange(filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, err
	}

	summaryModel, err := s.Repository.GetDashboardSummary(ctx, startDate, endDate, filter.OutletID)
	if err != nil {
		logger.WithError(err).Error("failed to get dashboard summary")
		return nil, err
	}

	orderStatusModels, err := s.Repository.GetOrderStatusRecap(ctx, startDate, endDate, filter.OutletID)
	if err != nil {
		logger.WithError(err).Error("failed to get order status recap")
		return nil, err
	}

	paymentMethodModels, err := s.Repository.GetPaymentMethodRecap(ctx, startDate, endDate, filter.OutletID)
	if err != nil {
		logger.WithError(err).Error("failed to get payment method recap")
		return nil, err
	}

	salesChartModels, err := s.Repository.GetSalesChart(ctx, startDate, endDate, filter.OutletID)
	if err != nil {
		logger.WithError(err).Error("failed to get sales chart")
		return nil, err
	}

	return s.mapToDashboardResponse(summaryModel, orderStatusModels, paymentMethodModels, salesChartModels), nil
}

func (s *ReportServiceImpl) GetFinancialReport(ctx context.Context, filter dto.FinancialReportFilter) (*dto.FinancialReportResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx)

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists || tenantID == uuid.Nil {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	if filter.OutletID != nil {
		if err := s.validateOutletAccess(ctx, *filter.OutletID); err != nil {
			return nil, err
		}
	}

	startDate, endDate, err := resolveDateRange(filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, err
	}

	setFilterDefault(&filter)

	summaryModel, err := s.Repository.GetFinancialSummary(ctx, startDate, endDate, filter.OutletID)
	if err != nil {
		logger.WithError(err).Error("failed to get financial summary")
		return nil, err
	}

	offset := (filter.Page - 1) * filter.Limit
	paymentModels, total, err := s.Repository.GetPayments(ctx, filter.OutletID, startDate, endDate, offset, filter.Limit)
	if err != nil {
		logger.WithError(err).Error("failed to get payments")
		return nil, err
	}

	paymentItems := make([]dto.PaymentReportItem, len(paymentModels))
	for i, p := range paymentModels {
		paymentItems[i] = dto.PaymentReportItem{
			PaymentID:     p.ID,
			OrderID:       p.OrderID,
			OrderNumber:   p.Order.OrderNumber,
			CustomerName:  p.Order.Customer.Name,
			OutletName:    p.Outlet.Name,
			PaymentMethod: p.PaymentMethod,
			Amount:        p.Amount,
			Status:        p.Status,
			PaidAt:        p.PaidAt,
		}
	}

	var totalPages int
	if total == 0 {
		totalPages = 0
	} else {
		totalPages = int(math.Ceil(float64(total) / float64(filter.Limit)))
	}

	return &dto.FinancialReportResponse{
		Data: dto.FinancialReportData{
			Summary: dto.FinancialSummary{
				TotalSales:            summaryModel.TotalSales,
				TotalDiscount:         summaryModel.TotalDiscount,
				TotalPaymentsReceived: summaryModel.TotalPaymentsReceived,
				TotalOrders:           summaryModel.TotalOrders,
			},
			Payments: paymentItems,
		},
		Pagination: dto.ReportPagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ReportServiceImpl) mapToDashboardResponse(
	summary models.DashboardSummary,
	orderStatuses []models.OrderStatusRecap,
	paymentMethods []models.PaymentMethodRecap,
	salesChart []models.SalesChartItem,
) *dto.DashboardResponse {
	statusRecaps := make([]dto.OrderStatusRecap, len(orderStatuses))
	for i, os := range orderStatuses {
		statusRecaps[i] = dto.OrderStatusRecap{
			Status: os.Status,
			Count:  os.Count,
		}
	}

	pmRecaps := make([]dto.PaymentMethodRecap, len(paymentMethods))
	for i, pm := range paymentMethods {
		pmRecaps[i] = dto.PaymentMethodRecap{
			Method: pm.Method,
			Total:  pm.Total,
			Count:  pm.Count,
		}
	}

	chartItems := make([]dto.SalesChartItem, len(salesChart))
	for i, sc := range salesChart {
		chartItems[i] = dto.SalesChartItem{
			Date:                  sc.Date,
			TotalSales:            sc.TotalSales,
			TotalPaymentsReceived: sc.TotalPaymentsReceived,
		}
	}

	return &dto.DashboardResponse{
		Summary: dto.DashboardSummary{
			TotalSales:            summary.TotalSales,
			TotalPaymentsReceived: summary.TotalPaymentsReceived,
			TotalOrders:           summary.TotalOrders,
			ActiveOrders:          summary.ActiveOrders,
		},
		OrderStatus:    statusRecaps,
		PaymentMethods: pmRecaps,
		SalesChart:     chartItems,
	}
}

func (s *ReportServiceImpl) validateOutletAccess(ctx context.Context, outletID uuid.UUID) error {
	role, _ := appcontext.RoleFromContext(ctx)
	if role == "ADMIN" {
		return nil
	}

	userID, ok := appcontext.UserIDFromContext(ctx)
	if !ok {
		return apperror.UnauthorizedError{Msg: "user_id not found"}
	}

	outletUser, err := s.OutletUserRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if outletUser.OutletID == outletID {
		return nil
	}
	return apperror.NotFoundError{Msg: "outlet not found or not accessible"}
}

func resolveDateRange(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	loc := defaultLocation
	var startDate, endDate time.Time
	var err error

	if startDateStr == "" && endDateStr == "" {
		now := time.Now().In(loc)
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, loc)
		return startDate, endDate, nil
	}

	if startDateStr != "" {
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, loc)
		if err != nil {
			return time.Time{}, time.Time{}, apperror.ValidationError{Msg: "invalid start_date format, use YYYY-MM-DD"}
		}
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	} else {
		startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, loc)
	}

	if endDateStr != "" {
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, loc)
		if err != nil {
			return time.Time{}, time.Time{}, apperror.ValidationError{Msg: "invalid end_date format, use YYYY-MM-DD"}
		}
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)
	} else {
		now := time.Now().In(loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	}

	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, apperror.ValidationError{Msg: "start_date cannot be after end_date"}
	}

	return startDate, endDate, nil
}

func setFilterDefault(f *dto.FinancialReportFilter) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 10
	} else if f.Limit > 100 {
		f.Limit = 100
	}
}
