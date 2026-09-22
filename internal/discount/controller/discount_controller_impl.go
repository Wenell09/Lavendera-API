package controller

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/Wenell09/lavendera-api/internal/discount/dto"
	"github.com/Wenell09/lavendera-api/internal/discount/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DiscountControllerImpl struct {
	Service   service.DiscountService
	Validator *validator.Validate
}

func NewDiscountController(service service.DiscountService, validator *validator.Validate) DiscountController {
	return &DiscountControllerImpl{Service: service, Validator: validator}
}

// Create implements [DiscountController].
func (d *DiscountControllerImpl) Create(c *gin.Context) {
	request := dto.CreateDiscountRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := d.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := d.Service.Create(c.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "discount created successfully", result))
}

// Delete implements [DiscountController].
func (d *DiscountControllerImpl) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid discount id"})
		return
	}
	if err := d.Service.Delete(c.Request.Context(), id); err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "discount deleted successfully", nil))
}

// FindAll implements [DiscountController].
func (d *DiscountControllerImpl) FindAll(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	filter := dto.DiscountFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	}
	result, err := d.Service.FindAll(c.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "discounts retrieved successfully", result))
}

// FindByID implements [DiscountController].
func (d *DiscountControllerImpl) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid discount id"})
		return
	}
	result, err := d.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "discount retrieved successfully", result))
}

// Update implements [DiscountController].
func (d *DiscountControllerImpl) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid discount id"})
		return
	}
	var request dto.UpdateDiscountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := d.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := d.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "discount updated successfully", result))
}
