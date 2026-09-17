package apperror

import (
	"errors"
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func NewHandleError(c *gin.Context, err error) {
	var appErr AppError
	if errors.As(err, &appErr) {
		c.JSON(
			appErr.StatusCode(),
			response.NewResponseError(
				appErr.StatusCode(),
				appErr.ResponseMessage(),
				appErr.ErrorData(),
			),
		)
		return
	}
	c.JSON(
		http.StatusInternalServerError,
		response.NewResponseError(
			http.StatusInternalServerError,
			"Internal Server Error",
			err.Error(),
		),
	)
}
