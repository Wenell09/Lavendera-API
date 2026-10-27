package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/report/dto"
	"github.com/Wenell09/lavendera-api/internal/report/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ReportControllerImpl struct {
	Service   service.ReportService
	Validator *validator.Validate
}

func NewReportController(service service.ReportService, validator *validator.Validate) ReportController {
	return &ReportControllerImpl{
		Service:   service,
		Validator: validator,
	}
}

func (c *ReportControllerImpl) GetDashboard(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	var outletID *uuid.UUID
	if outletParam := ctx.Query("outlet_id"); outletParam != "" {
		id, err := uuid.Parse(outletParam)
		if err != nil {
			apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid outlet_id format"})
			return
		}
		outletID = &id
	}

	filter := dto.DashboardFilter{
		StartDate: startDate,
		EndDate:   endDate,
		OutletID:  outletID,
	}

	result, err := c.Service.GetDashboard(ctx.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "dashboard retrieved successfully", result, response.ResponseMeta{}))
}

func (c *ReportControllerImpl) GetFinancialReport(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil || p <= 0 {
			apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid page format"})
			return
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil || l <= 0 {
			apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid limit format"})
			return
		}
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var outletID *uuid.UUID
	if outletParam := ctx.Query("outlet_id"); outletParam != "" {
		id, err := uuid.Parse(outletParam)
		if err != nil {
			apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid outlet_id format"})
			return
		}
		outletID = &id
	}

	filter := dto.FinancialReportFilter{
		StartDate: startDate,
		EndDate:   endDate,
		OutletID:  outletID,
		Page:      page,
		Limit:     limit,
	}

	result, err := c.Service.GetFinancialReport(ctx.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "financial report retrieved successfully", result.Data, response.ResponseMeta{Pagination: result.Pagination}))
}
