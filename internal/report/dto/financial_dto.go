package dto

import (
	"time"

	"github.com/google/uuid"
)

type FinancialReportFilter struct {
	StartDate string     `form:"start_date"`
	EndDate   string     `form:"end_date"`
	OutletID  *uuid.UUID `form:"outlet_id"`
	Page      int        `form:"page"`
	Limit     int        `form:"limit"`
}

type FinancialSummary struct {
	TotalSales            int64 `json:"total_sales"`
	TotalDiscount         int64 `json:"total_discount"`
	TotalPaymentsReceived int64 `json:"total_payments_received"`
	TotalOrders           int64 `json:"total_orders"`
}

type PaymentReportItem struct {
	PaymentID     uuid.UUID  `json:"payment_id"`
	OrderID       uuid.UUID  `json:"order_id"`
	OrderNumber   string     `json:"order_number"`
	CustomerName  string     `json:"customer_name"`
	OutletName    string     `json:"outlet_name"`
	PaymentMethod string     `json:"payment_method"`
	Amount        int64      `json:"amount"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
}

type FinancialReportData struct {
	Summary  FinancialSummary    `json:"summary"`
	Payments []PaymentReportItem `json:"payments"`
}

type ReportPagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type FinancialReportResponse struct {
	Data       FinancialReportData `json:"data"`
	Pagination ReportPagination    `json:"pagination"`
}
