package apperror

import "net/http"

type ConflictError struct {
	Msg string
}

func (e ConflictError) Error() string {
	return e.Msg
}

func (e ConflictError) StatusCode() int {
	return http.StatusConflict
}

func (e ConflictError) ResponseMessage() string {
	return "Conflict"
}

func (e ConflictError) ErrorData() interface{} {
	return e.Msg
}

type UnauthorizedError struct {
	Msg string
}

func (e UnauthorizedError) Error() string {
	return e.Msg
}

func (e UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

func (e UnauthorizedError) ResponseMessage() string {
	return "Unauthorized"
}

func (e UnauthorizedError) ErrorData() interface{} {
	return e.Msg
}
