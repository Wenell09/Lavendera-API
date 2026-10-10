package dto

import "github.com/google/uuid"

type OrderFilter struct {
	Search     string     `form:"search"`
	Status     string     `form:"status"`
	OutletID   *uuid.UUID `form:"outlet_id"`
	CustomerID *uuid.UUID `form:"customer_id"`
	StartDate  string     `form:"start_date"`
	EndDate    string     `form:"end_date"`
	Page       int        `form:"page"`
	Limit      int        `form:"limit"`
}

func (f *OrderFilter) SetDefault() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 10
	} else if f.Limit > 100 {
		f.Limit = 100
	}
}
