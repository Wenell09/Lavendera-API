package utils

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers/v2"
)

func IsValidIndonesianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if !strings.HasPrefix(phone, "62") {
		return false
	}
	if len(phone) < 12 || len(phone) > 15 {
		return false
	}
	if _, err := strconv.ParseUint(phone, 10, 64); err != nil {
		return false
	}
	num, err := phonenumbers.Parse("+"+phone, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumberForRegion(num, "ID")
}
