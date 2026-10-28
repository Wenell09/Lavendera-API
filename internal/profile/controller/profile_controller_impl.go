package controller

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/profile/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type ProfileControllerImpl struct {
	Service service.ProfileService
}

func NewProfileController(service service.ProfileService) ProfileController {
	return &ProfileControllerImpl{
		Service: service,
	}
}

// GetProfile implements [ProfileController].
func (p *ProfileControllerImpl) GetProfile(c *gin.Context) {
	result, err := p.Service.GetProfile(c.Request.Context())
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "profile retrieved successfully", result, response.ResponseMeta{}))
}
