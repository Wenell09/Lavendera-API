package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers/v2"
)

func IsValidIndonesianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	num, err := phonenumbers.Parse(phone, "ID")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumberForRegion(num, "ID")
}
