package controller

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/auth/dto"
	"github.com/Wenell09/lavendera-api/internal/auth/service"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthControllerImpl struct {
	Service   service.AuthService
	Validator *validator.Validate
	Logger    *logrus.Logger
}

func NewAuthController(
	service service.AuthService,
	validator *validator.Validate,
	logger *logrus.Logger,
) AuthController {
	return &AuthControllerImpl{
		Service:   service,
		Validator: validator,
		Logger:    logger,
	}
}

// Login implements [AuthController].
func (a *AuthControllerImpl) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.Logger.WithError(err).
			Warn("login failed: invalid JSON")
		c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				http.StatusBadRequest,
				"Invalid request body",
				err.Error(),
			),
		)
		return
	}
	// DTO validation
	if err := a.Validator.Struct(req); err != nil {
		validationErr, ok :=
			err.(validator.ValidationErrors)
		if ok {
			apperror.NewHandleError(
				c,
				apperror.NewFieldError(validationErr),
			)
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	// Service
	result, err := a.Service.Login(
		c.Request.Context(),
		req,
	)
	if err != nil {
		a.Logger.WithError(err).
			Warn("login failed")
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(
		http.StatusOK,
		response.NewResponseSuccess(
			http.StatusOK,
			"Login successful",
			result,
		),
	)
}

// Register implements [AuthController].
func (a *AuthControllerImpl) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.Logger.WithError(err).
			Warn("register failed: invalid JSON")
		c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				http.StatusBadRequest,
				"Invalid request body",
				err.Error(),
			),
		)
		return
	}
	// DTO validation
	if err := a.Validator.Struct(req); err != nil {
		validationErr, ok := err.(validator.ValidationErrors)
		if ok {
			apperror.NewHandleError(
				c,
				apperror.NewFieldError(validationErr),
			)
			return
		}
		apperror.NewHandleError(c, err)
		return
	}
	// Service
	result, err := a.Service.Register(
		c.Request.Context(),
		req,
	)
	if err != nil {
		a.Logger.WithError(err).
			Warn("register failed")
		apperror.NewHandleError(c, err)
		return
	}
	c.JSON(
		http.StatusCreated,
		response.NewResponseSuccess(
			http.StatusCreated,
			"Registration successful",
			result,
		),
	)
}
