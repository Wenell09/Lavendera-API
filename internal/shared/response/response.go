package response

type ResponseSuccess struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
	Data    interface{}  `json:"data,omitempty"`
	Meta    ResponseMeta `json:"meta"`
}

type ResponseError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error,omitempty"`
}

type ResponseMeta struct {
	Pagination interface{} `json:"pagination,omitempty"`
}

func NewResponseSuccess(status int, message string, data interface{}, meta ResponseMeta) ResponseSuccess {
	return ResponseSuccess{
		Status:  status,
		Message: message,
		Success: status >= 200 && status < 300,
		Data:    data,
		Meta:    meta,
	}
}

func NewResponseError(status int, message string, error interface{}) ResponseError {
	return ResponseError{
		Status:  status,
		Message: message,
		Success: status >= 200 && status < 300,
		Error:   error,
	}
}
