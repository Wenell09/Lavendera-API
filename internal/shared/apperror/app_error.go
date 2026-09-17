package apperror

type AppError interface {
	error
	StatusCode() int
	ResponseMessage() string
	ErrorData() interface{}
}
