package apperror

import "net/http"

type NotFoundError struct {
	Msg string
}

func (n NotFoundError) Error() string {
	return n.Msg
}

func (n NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func (n NotFoundError) ResponseMessage() string {
	return "Not Found"
}

func (n NotFoundError) ErrorData() interface{} {
	return n.Msg
}
