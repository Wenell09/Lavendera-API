package service

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/report/dto"
)

type ReportService interface {
	GetDashboard(ctx context.Context, filter dto.DashboardFilter) (*dto.DashboardResponse, error)
	GetFinancialReport(ctx context.Context, filter dto.FinancialReportFilter) (*dto.FinancialReportResponse, error)
}
