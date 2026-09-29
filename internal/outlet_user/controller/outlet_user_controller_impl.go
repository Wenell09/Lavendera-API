package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/outlet_user/dto"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OutletUserControllerImpl struct {
	Service   service.OutletUserService
	Validator *validator.Validate
}

func NewOutletUserController(service service.OutletUserService, validator *validator.Validate) OutletUserController {
	return &OutletUserControllerImpl{Service: service, Validator: validator}
}

func (u *OutletUserControllerImpl) Assign(c *gin.Context) {
	var request dto.CreateOutletUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}

	if err := u.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}

	result, err := u.Service.Assign(c.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "user assigned to outlet successfully", result, response.ResponseMeta{}))
}

func (u *OutletUserControllerImpl) Unassign(c *gin.Context) {
	outletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid user id"})
		return
	}

	if err := u.Service.Unassign(c.Request.Context(), outletID, userID); err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "user unassigned from outlet successfully", nil, response.ResponseMeta{}))
}

func (u *OutletUserControllerImpl) FindByOutletID(c *gin.Context) {
	outletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	filter := dto.OutletUserFilter{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	result, err := u.Service.FindByOutletID(c.Request.Context(), outletID, filter)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet users retrieved successfully", result.Data, response.ResponseMeta{
		Pagination: result.Pagination,
	}))
}

