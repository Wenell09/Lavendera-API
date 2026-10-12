package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/order/dto"
	"github.com/Wenell09/lavendera-api/internal/order/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OrderControllerImpl struct {
	Service   service.OrderService
	Validator *validator.Validate
}

func NewOrderController(service service.OrderService, validator *validator.Validate) OrderController {
	return &OrderControllerImpl{
		Service:   service,
		Validator: validator,
	}
}

// CreateCashierOrder implements [OrderController].
func (c *OrderControllerImpl) CreateCashierOrder(ctx *gin.Context) {
	var request dto.CreateCashierOrderRequest
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

	result, err := c.Service.CreateCashierOrder(ctx.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "order created successfully", result, response.ResponseMeta{}))
}

// CreatePublicOrder implements [OrderController].
func (c *OrderControllerImpl) CreatePublicOrder(ctx *gin.Context) {
	var request dto.CreatePublicOrderRequest
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

	result, err := c.Service.CreatePublicOrder(ctx.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "public order created successfully", result, response.ResponseMeta{}))
}

// FindAll implements [OrderController].
func (c *OrderControllerImpl) FindAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	status := ctx.Query("status")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	var outletID *uuid.UUID
	if outletParam := ctx.Query("outlet_id"); outletParam != "" {
		if id, err := uuid.Parse(outletParam); err == nil {
			outletID = &id
		}
	}

	var customerID *uuid.UUID
	if custParam := ctx.Query("customer_id"); custParam != "" {
		if id, err := uuid.Parse(custParam); err == nil {
			customerID = &id
		}
	}

	filter := dto.OrderFilter{
		Search:     search,
		Status:     status,
		OutletID:   outletID,
		CustomerID: customerID,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       page,
		Limit:      limit,
	}
	filter.SetDefault()

	result, err := c.Service.FindAll(ctx.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "orders retrieved successfully", result.Data, response.ResponseMeta{Pagination: result.Pagination}))
}

// FindByID implements [OrderController].
func (c *OrderControllerImpl) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid order id"})
		return
	}

	result, err := c.Service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "order retrieved successfully", result, response.ResponseMeta{}))
}

// TrackOrderByToken implements [OrderController].
func (c *OrderControllerImpl) TrackOrderByToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "tracking token is required"})
		return
	}

	result, err := c.Service.FindByTrackingToken(ctx.Request.Context(), token)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "order tracked successfully", result, response.ResponseMeta{}))
}

// UpdateStatus implements [OrderController].
func (c *OrderControllerImpl) UpdateStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid order id"})
		return
	}

	var request dto.UpdateOrderStatusRequest
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

	result, err := c.Service.UpdateStatus(ctx.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "order status updated successfully", result, response.ResponseMeta{}))
}

// CancelOrder implements [OrderController].
func (c *OrderControllerImpl) CancelOrder(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid order id"})
		return
	}

	result, err := c.Service.CancelOrder(ctx.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "order canceled successfully", result, response.ResponseMeta{}))
}

// ConfirmPayment implements [OrderController].
func (c *OrderControllerImpl) ConfirmPayment(ctx *gin.Context) {
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid order id"})
		return
	}

	paymentID, err := uuid.Parse(ctx.Param("payment_id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid payment id"})
		return
	}

	result, err := c.Service.ConfirmPayment(ctx.Request.Context(), orderID, paymentID)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "payment confirmed successfully", result, response.ResponseMeta{}))
}

// InputCashPayment implements [OrderController].
func (c *OrderControllerImpl) InputCashPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid order id"})
		return
	}

	var request dto.UpdateOrderPaidAmountRequest
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

	result, err := c.Service.InputCashPayment(ctx.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "cash payment input successfully", result, response.ResponseMeta{}))
}
