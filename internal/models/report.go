package models

type DashboardSummary struct {
	TotalSales            int64
	TotalPaymentsReceived int64
	TotalOrders           int64
	ActiveOrders          int64
}

type OrderStatusRecap struct {
	Status string
	Count  int64
}

type PaymentMethodRecap struct {
	Method string
	Total  int64
	Count  int64
}

type SalesChartItem struct {
	Date                  string
	TotalSales            int64
	TotalPaymentsReceived int64
}

type FinancialSummary struct {
	TotalSales            int64
	TotalDiscount         int64
	TotalPaymentsReceived int64
	TotalOrders           int64
}
