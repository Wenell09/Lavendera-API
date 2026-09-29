package controller

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/Wenell09/lavendera-api/internal/customer/dto"
	"github.com/Wenell09/lavendera-api/internal/customer/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerControllerImpl struct {
	Service   service.CustomerService
	Validator *validator.Validate
}

func NewCustomerController(service service.CustomerService, validator *validator.Validate) CustomerController {
	return &CustomerControllerImpl{Service: service, Validator: validator}
}

// Create implements [CustomerController].
func (c *CustomerControllerImpl) Create(ctx *gin.Context) {
	request := dto.CreateCustomerRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := c.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(ctx, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(ctx, err)
		return
	}
	result, err := c.Service.Create(ctx.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "customer created successfully", result, response.ResponseMeta{}))
}

// Delete implements [CustomerController].
func (c *CustomerControllerImpl) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid customer id"})
		return
	}
	if err := c.Service.Delete(ctx.Request.Context(), id); err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "customer deleted successfully", nil, response.ResponseMeta{}))
}

// FindAll implements [CustomerController].
func (c *CustomerControllerImpl) FindAll(ctx *gin.Context) {
	search := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	filter := dto.CustomerFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	}
	result, err := c.Service.FindAll(ctx.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "customers retrieved successfully", result.Data, response.ResponseMeta{Pagination: result.Pagination}))
}

// FindByID implements [CustomerController].
func (c *CustomerControllerImpl) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid customer id"})
		return
	}
	result, err := c.Service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "customer retrieved successfully", result, response.ResponseMeta{}))
}

// Update implements [CustomerController].
func (c *CustomerControllerImpl) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid customer id"})
		return
	}
	var request dto.UpdateCustomerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := c.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(ctx, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(ctx, err)
		return
	}
	result, err := c.Service.Update(ctx.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "customer updated successfully", result, response.ResponseMeta{}))
}
