package controller

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/service_category/dto"
	"github.com/Wenell09/lavendera-api/internal/service_category/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ServiceCategoryControllerImpl struct {
	Service   service.ServiceCategoryService
	Validator *validator.Validate
}

func NewServiceCategoryController(service service.ServiceCategoryService, validator *validator.Validate) ServiceCategoryController {
	return &ServiceCategoryControllerImpl{Service: service, Validator: validator}
}

func (s *ServiceCategoryControllerImpl) Create(c *gin.Context) {
	var req dto.CreateServiceCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := s.Validator.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := s.Service.Create(c.Request.Context(), req)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "service category created successfully", result, response.ResponseMeta{}))
}

func (s *ServiceCategoryControllerImpl) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid service id"})
		return
	}
	if err := s.Service.Delete(c.Request.Context(), id); err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service category deleted successfully", nil, response.ResponseMeta{}))
}

func (s *ServiceCategoryControllerImpl) FindAll(c *gin.Context) {
	result, err := s.Service.FindAll(c.Request.Context())
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service categories retrieved successfully", result, response.ResponseMeta{}))
}

func (s *ServiceCategoryControllerImpl) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid service id"})
		return
	}
	result, err := s.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service category retrieved successfully", result, response.ResponseMeta{}))
}

func (s *ServiceCategoryControllerImpl) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid service id"})
		return
	}
	var req dto.UpdateServiceCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := s.Validator.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := s.Service.Update(c.Request.Context(), id, req)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service category updated successfully", result, response.ResponseMeta{}))
}
