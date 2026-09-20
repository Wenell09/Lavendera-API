package dto

import "github.com/google/uuid"

type ServiceFilter struct {
	Search     string     `json:"search"`
	CategoryID *uuid.UUID `json:"category_id"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
}
