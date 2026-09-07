package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHandleError(c *gin.Context, err error) {
	if appErr, ok := err.(AppError); ok {
		c.JSON(appErr.StatusCode(), NewResponseError(
			appErr.StatusCode(),
			appErr.ResponseMessage(),
			appErr.ErrorData(),
		))
		return
	}

	c.JSON(http.StatusInternalServerError, NewResponseError(
		http.StatusInternalServerError,
		"Internal Server Error",
		err.Error(),
	))
}
