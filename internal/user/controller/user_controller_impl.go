package controller

import (
	"net/http"
	"strconv"

	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/Wenell09/lavendera-api/internal/user/dto"
	"github.com/Wenell09/lavendera-api/internal/user/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserControllerImpl struct {
	Service   service.UserService
	Validator *validator.Validate
}

func NewUserController(service service.UserService, validator *validator.Validate) UserController {
	return &UserControllerImpl{Service: service, Validator: validator}
}

// Create implements [UserController].
func (u *UserControllerImpl) Create(c *gin.Context) {
	request := dto.CreateUserRequest{}
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

	result, err := u.Service.Create(c.Request.Context(), request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.NewResponseSuccess(http.StatusCreated, "user created successfully", result, response.ResponseMeta{}))
}

// Delete implements [UserController].
func (u *UserControllerImpl) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid user id"})
		return
	}

	if err := u.Service.Delete(c.Request.Context(), id); err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "user deleted successfully", nil, response.ResponseMeta{}))
}

// FindAll implements [UserController].
func (u *UserControllerImpl) FindAll(c *gin.Context) {
	search := c.Query("search")
	role := c.Query("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := dto.UserFilter{
		Search: search,
		Role:   role,
		Page:   page,
		Limit:  limit,
	}

	result, err := u.Service.FindAll(c.Request.Context(), filter)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "users retrieved successfully", result.Data, response.ResponseMeta{Pagination: result.Pagination}))
}

// FindByID implements [UserController].
func (u *UserControllerImpl) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid user id"})
		return
	}

	result, err := u.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "user retrieved successfully", result, response.ResponseMeta{}))
}

// Update implements [UserController].
func (u *UserControllerImpl) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperror.NewHandleError(c, apperror.ValidationError{Msg: "invalid user id"})
		return
	}

	var request dto.UpdateUserRequest
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

	result, err := u.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		apperror.NewHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(http.StatusOK, "user updated successfully", result, response.ResponseMeta{}))
}
