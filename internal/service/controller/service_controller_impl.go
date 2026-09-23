package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/service/dto"
	"github.com/Wenell09/lavendera-api/internal/service/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ServiceControllerImpl struct {
	Service   service.Service
	Validator *validator.Validate
}

func NewServiceController(service service.Service, validator *validator.Validate) ServiceController {
	return &ServiceControllerImpl{Service: service, Validator: validator}
}

func (s *ServiceControllerImpl) Create(c *gin.Context) {
	var request dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := s.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := s.Service.Create(c.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "service created successfully", result, response.ResponseMeta{}))
}

func (s *ServiceControllerImpl) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid service id"})
		return
	}
	if err := s.Service.Delete(c.Request.Context(), id); err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service deleted successfully", nil, response.ResponseMeta{}))
}

func (s *ServiceControllerImpl) FindAll(c *gin.Context) {
	search := c.Query("search")
	categoryIDParam := c.Query("category_id")
	var categoryID *uuid.UUID
	if categoryIDParam != "" {
		id, err := uuid.Parse(categoryIDParam)
		if err != nil {
			apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid category_id"})
			return
		}
		categoryID = &id
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	filter := dto.ServiceFilter{
		Search:     search,
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
	}
	result, err := s.Service.FindAll(c.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "services retrieved successfully", result.Data, response.ResponseMeta{Pagination: result.Pagination}))
}

func (s *ServiceControllerImpl) FindByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service retrieved successfully", result, response.ResponseMeta{}))
}

func (s *ServiceControllerImpl) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid service id"})
		return
	}
	var request dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid request body"})
		return
	}
	if err := s.Validator.Struct(request); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			apperror.NewHandleError(c, apperror.NewFieldError(ve))
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	result, err := s.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "service updated successfully", result, response.ResponseMeta{}))
}
