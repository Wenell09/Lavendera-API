package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/dto"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OutletPaymentMethodControllerImpl struct {
	Service   service.OutletPaymentMethodService
	Validator *validator.Validate
}

func NewOutletPaymentMethodController(s service.OutletPaymentMethodService, v *validator.Validate) OutletPaymentMethodController {
	return &OutletPaymentMethodControllerImpl{
		Service:   s,
		Validator: v,
	}
}

func (ctrl *OutletPaymentMethodControllerImpl) Create(c *gin.Context) {
	outletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}

	var req dto.CreateOutletPaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}

	req.OutletID = outletID
	if err := ctrl.Validator.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}

	res, err := ctrl.Service.Create(c.Request.Context(), req)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "outlet payment method created successfully", res, response.ResponseMeta{}))
}

func (ctrl *OutletPaymentMethodControllerImpl) FindAll(c *gin.Context) {
	outletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid outlet id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	filter := dto.OutletPaymentMethodFilter{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	res, err := ctrl.Service.FindAll(c.Request.Context(), outletID, filter)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet payment methods retrieved successfully", res.Data, response.ResponseMeta{
		Pagination: res.Pagination,
	}))
}

func (ctrl *OutletPaymentMethodControllerImpl) FindByID(c *gin.Context) {
	paymentMethodID, err := uuid.Parse(c.Param("payment_method_id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid payment method id"})
		return
	}

	res, err := ctrl.Service.FindByID(c.Request.Context(), paymentMethodID)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet payment method retrieved successfully", res, response.ResponseMeta{}))
}

func (ctrl *OutletPaymentMethodControllerImpl) Update(c *gin.Context) {
	paymentMethodID, err := uuid.Parse(c.Param("payment_method_id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid payment method id"})
		return
	}

	var req dto.UpdateOutletPaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}

	if err := ctrl.Validator.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}

	res, err := ctrl.Service.Update(c.Request.Context(), paymentMethodID, req)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet payment method updated successfully", res, response.ResponseMeta{}))
}

func (ctrl *OutletPaymentMethodControllerImpl) Delete(c *gin.Context) {
	paymentMethodID, err := uuid.Parse(c.Param("payment_method_id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid payment method id"})
		return
	}

	err = ctrl.Service.Delete(c.Request.Context(), paymentMethodID)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "outlet payment method deleted successfully", nil, response.ResponseMeta{}))
}
