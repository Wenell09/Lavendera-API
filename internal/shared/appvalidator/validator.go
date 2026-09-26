package appvalidator

import (
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("id_phone", utils.IsValidIndonesianPhone)
	return v
}
