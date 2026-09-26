package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/outlet/dto"
	"github.com/Wenell09/lavendera-api/internal/outlet/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OutletControllerImpl struct {
	Service   service.OutletService
	Validator *validator.Validate
}

func NewOutletController(service service.OutletService, validator *validator.Validate) OutletController {
	return &OutletControllerImpl{Service: service, Validator: validator}
}

// Create implements [OutletController].
func (o *OutletControllerImpl) Create(c *gin.Context) {
	request := dto.CreateOutletRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := o.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := o.Service.Create(c.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "outlet created successfully", result, response.ResponseMeta{}))
}

// Delete implements [OutletController].
func (o *OutletControllerImpl) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}
	if err := o.Service.Delete(c.Request.Context(), id); err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet deleted successfully", nil, response.ResponseMeta{}))
}

// FindAll implements [OutletController].
func (o *OutletControllerImpl) FindAll(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	filter := dto.OutletFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	}
	result, err := o.Service.FindAll(c.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlets retrieved successfully", result.Data, response.ResponseMeta{Pagination: result.Pagination}))
}

// FindByID implements [OutletController].
func (o *OutletControllerImpl) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}
	result, err := o.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet retrieved successfully", result, response.ResponseMeta{}))
}

// Update implements [OutletController].
func (o *OutletControllerImpl) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}
	var request dto.UpdateOutletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := o.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := o.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet updated successfully", result, response.ResponseMeta{}))
}
