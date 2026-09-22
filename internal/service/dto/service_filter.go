package dto

import "github.com/google/uuid"

type ServiceFilter struct {
	Search     string     `json:"search"`
	CategoryID *uuid.UUID `json:"category_id"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
}

func (f *ServiceFilter) SetDefault() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
}
