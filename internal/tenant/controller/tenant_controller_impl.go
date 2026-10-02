package controller

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/Wenell09/lavendera-api/internal/tenant/dto"
	"github.com/Wenell09/lavendera-api/internal/tenant/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TenantControllerImpl struct {
	Service   service.TenantService
	Validator *validator.Validate
}

func NewTenantController(service service.TenantService, validator *validator.Validate) TenantController {
	return &TenantControllerImpl{Service: service, Validator: validator}
}

// FindByID implements [TenantController]
func (c *TenantControllerImpl) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid tenant id"})
		return
	}
	result, err := c.Service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "tenant retrieved successfully", result, response.ResponseMeta{}))
}

// Update implements [TenantController]
func (c *TenantControllerImpl) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid tenant id"})
		return
	}
	var req dto.UpdateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apperror.NewHandleError(ctx, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := c.Validator.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(ctx, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(ctx, err)
		return
	}
	result, err := c.Service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		apperror.NewHandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "tenant updated successfully", result, response.ResponseMeta{}))
}
