package dto

import (
	"github.com/google/uuid"
)

type DashboardFilter struct {
	StartDate string     `form:"start_date"`
	EndDate   string     `form:"end_date"`
	OutletID  *uuid.UUID `form:"outlet_id"`
}

type DashboardSummary struct {
	TotalSales            int64 `json:"total_sales"`
	TotalPaymentsReceived int64 `json:"total_payments_received"`
	TotalOrders           int64 `json:"total_orders"`
	ActiveOrders          int64 `json:"active_orders"`
}

type OrderStatusRecap struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type PaymentMethodRecap struct {
	Method string `json:"method"`
	Total  int64  `json:"total"`
	Count  int64  `json:"count"`
}

type SalesChartItem struct {
	Date                  string `json:"date"`
	TotalSales            int64  `json:"total_sales"`
	TotalPaymentsReceived int64  `json:"total_payments_received"`
}

type DashboardResponse struct {
	Summary        DashboardSummary     `json:"summary"`
	OrderStatus    []OrderStatusRecap   `json:"order_status"`
	PaymentMethods []PaymentMethodRecap `json:"payment_methods"`
	SalesChart     []SalesChartItem     `json:"sales_chart"`
}
