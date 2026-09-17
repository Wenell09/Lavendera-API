package apperror

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

// validation manual
type ValidationError struct {
	Msg string
}

func (v ValidationError) Error() string {
	return v.Msg
}

func (v ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

func (v ValidationError) ResponseMessage() string {
	return "Validation Error"
}

func (v ValidationError) ErrorData() interface{} {
	return v.Msg
}

// validation menggunakan validator
type FieldDetail struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type FieldError struct {
	Errors []FieldDetail `json:"errors"`
}

func (f FieldError) Error() string {
	return "validation error"
}

func (f FieldError) StatusCode() int {
	return http.StatusBadRequest
}

func (f FieldError) ResponseMessage() string {
	return "Validation Error"
}

func (f FieldError) ErrorData() interface{} {
	return f.Errors
}

func NewFieldError(validationErr validator.ValidationErrors) error {
	dataErr := []FieldDetail{}

	for _, fieldErr := range validationErr {
		dataErr = append(dataErr, FieldDetail{
			Field: fieldErr.Field(),
			Error: fieldErr.Tag(),
		})
	}

	return FieldError{
		Errors: dataErr,
	}
}
