package repository

import (
	"context"
	"time"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepositoryImpl struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{DB: db}
}

func (r *ReportRepositoryImpl) GetDashboardSummary(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.DashboardSummary, error) {
	var summary models.DashboardSummary

	var orderStats struct {
		TotalSales int64
	}
	query := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status != ?", "DIBATALKAN")

	if outletID != nil {
		query = query.Where("outlet_id = ?", *outletID)
	}

	if err := query.Select("COALESCE(SUM(total), 0) as total_sales").Scan(&orderStats).Error; err != nil {
		return summary, err
	}
	summary.TotalSales = orderStats.TotalSales

	queryTotalOrders := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	if outletID != nil {
		queryTotalOrders = queryTotalOrders.Where("outlet_id = ?", *outletID)
	}
	if err := queryTotalOrders.Count(&summary.TotalOrders).Error; err != nil {
		return summary, err
	}

	queryPayment := r.DB.WithContext(ctx).Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).
		Where("status = ?", "verified").
		Where("verified_at >= ? AND verified_at <= ?", startDate, endDate)

	if outletID != nil {
		queryPayment = queryPayment.Where("outlet_id = ?", *outletID)
	}

	if err := queryPayment.Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalPaymentsReceived).Error; err != nil {
		return summary, err
	}

	queryActive := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("status IN ?", []string{"MENUNGGU_KONFIRMASI", "DIPROSES"})

	if outletID != nil {
		queryActive = queryActive.Where("outlet_id = ?", *outletID)
	}

	if err := queryActive.Count(&summary.ActiveOrders).Error; err != nil {
		return summary, err
	}

	return summary, nil
}

func (r *ReportRepositoryImpl) GetOrderStatusRecap(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.OrderStatusRecap, error) {
	var recap []models.OrderStatusRecap

	statuses := []string{"MENUNGGU_KONFIRMASI", "DIPROSES", "SELESAI", "DIAMBIL", "DIBATALKAN"}

	query := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)

	if outletID != nil {
		query = query.Where("outlet_id = ?", *outletID)
	}

	var results []struct {
		Status string
		Count  int64
	}
	if err := query.Select("status, COUNT(id) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[string]int64)
	for _, res := range results {
		countMap[res.Status] = res.Count
	}

	for _, s := range statuses {
		recap = append(recap, models.OrderStatusRecap{
			Status: s,
			Count:  countMap[s],
		})
	}

	return recap, nil
}

func (r *ReportRepositoryImpl) GetPaymentMethodRecap(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.PaymentMethodRecap, error) {
	var recap []models.PaymentMethodRecap

	query := r.DB.WithContext(ctx).Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).
		Where("status = ?", "verified").
		Where("verified_at >= ? AND verified_at <= ?", startDate, endDate)

	if outletID != nil {
		query = query.Where("outlet_id = ?", *outletID)
	}

	if err := query.Select("payment_method as method, COALESCE(SUM(amount), 0) as total, COUNT(id) as count").
		Group("payment_method").Scan(&recap).Error; err != nil {
		return nil, err
	}

	return recap, nil
}

func (r *ReportRepositoryImpl) GetSalesChart(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) ([]models.SalesChartItem, error) {
	var salesResults []struct {
		Date  time.Time
		Total int64
	}
	salesQuery := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status != ?", "DIBATALKAN").
		Select("DATE(created_at) as date, SUM(total) as total").
		Group("DATE(created_at)").
		Order("date ASC")

	if outletID != nil {
		salesQuery = salesQuery.Where("outlet_id = ?", *outletID)
	}

	if err := salesQuery.Scan(&salesResults).Error; err != nil {
		return nil, err
	}

	var paymentResults []struct {
		Date  time.Time
		Total int64
	}
	paymentQuery := r.DB.WithContext(ctx).Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).
		Where("status = ?", "verified").
		Where("verified_at >= ? AND verified_at <= ?", startDate, endDate).
		Select("DATE(verified_at) as date, SUM(amount) as total").
		Group("DATE(verified_at)").
		Order("date ASC")

	if outletID != nil {
		paymentQuery = paymentQuery.Where("outlet_id = ?", *outletID)
	}

	if err := paymentQuery.Scan(&paymentResults).Error; err != nil {
		return nil, err
	}

	salesMap := make(map[string]int64)
	for _, res := range salesResults {
		salesMap[res.Date.Format("2006-01-02")] = res.Total
	}

	paymentMap := make(map[string]int64)
	for _, res := range paymentResults {
		paymentMap[res.Date.Format("2006-01-02")] = res.Total
	}

	var chart []models.SalesChartItem
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		chart = append(chart, models.SalesChartItem{
			Date:                  dateStr,
			TotalSales:            salesMap[dateStr],
			TotalPaymentsReceived: paymentMap[dateStr],
		})
	}

	return chart, nil
}

func (r *ReportRepositoryImpl) GetFinancialSummary(ctx context.Context, startDate, endDate time.Time, outletID *uuid.UUID) (models.FinancialSummary, error) {
	var summary models.FinancialSummary

	var orderStats struct {
		TotalSales    int64
		TotalDiscount int64
	}
	query := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Where("status != ?", "DIBATALKAN")

	if outletID != nil {
		query = query.Where("outlet_id = ?", *outletID)
	}

	if err := query.Select("COALESCE(SUM(total), 0) as total_sales, COALESCE(SUM(discount_amount), 0) as total_discount").Scan(&orderStats).Error; err != nil {
		return summary, err
	}
	summary.TotalSales = orderStats.TotalSales
	summary.TotalDiscount = orderStats.TotalDiscount

	queryTotalOrders := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	if outletID != nil {
		queryTotalOrders = queryTotalOrders.Where("outlet_id = ?", *outletID)
	}
	if err := queryTotalOrders.Count(&summary.TotalOrders).Error; err != nil {
		return summary, err
	}

	queryPayment := r.DB.WithContext(ctx).Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).
		Where("status = ?", "verified").
		Where("verified_at >= ? AND verified_at <= ?", startDate, endDate)

	if outletID != nil {
		queryPayment = queryPayment.Where("outlet_id = ?", *outletID)
	}

	if err := queryPayment.Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalPaymentsReceived).Error; err != nil {
		return summary, err
	}

	return summary, nil
}

func (r *ReportRepositoryImpl) GetPayments(ctx context.Context, outletID *uuid.UUID, startDate, endDate time.Time, offset, limit int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	dbQuery := r.DB.WithContext(ctx).Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).
		Where("payments.verified_at >= ? AND payments.verified_at <= ?", startDate, endDate).
		Where("payments.status = ?", "verified")

	if outletID != nil {
		dbQuery = dbQuery.Where("payments.outlet_id = ?", *outletID)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.
		Preload("Outlet").
		Preload("Order").
		Preload("Order.Customer").
		Order("payments.verified_at DESC, payments.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
