package utils

type ResponseSuccess struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error,omitempty"`
}

func NewResponseSuccess(status int, message string, data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Status:  status,
		Message: message,
		Success: status >= 200 && status < 300,
		Data:    data,
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
