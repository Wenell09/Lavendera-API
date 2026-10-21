package apperror

import (
	"net/http"
	"strings"

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
type FieldError struct {
	Msg string
}

func (f FieldError) Error() string {
	return f.Msg
}

func (f FieldError) StatusCode() int {
	return http.StatusBadRequest
}

func (f FieldError) ResponseMessage() string {
	return "Validation Error"
}

func (f FieldError) ErrorData() interface{} {
	return f.Msg
}

func NewFieldError(validationErr validator.ValidationErrors) error {
	messages := make([]string, 0, len(validationErr))
	for _, fieldErr := range validationErr {
		messages = append(
			messages,
			strings.ToLower(fieldErr.Field())+" "+fieldErr.Tag(),
		)
	}
	return FieldError{
		Msg: strings.Join(messages, ", "),
	}
}
