package app

import (
	"net/http"

	"github.com/Wenell09/lavendera-api/internal/auth/controller"
	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/Wenell09/lavendera-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewRouter(
	authController controller.AuthController,
	jwtConfig config.JWTConfig,
	db *gorm.DB,
	logger *logrus.Logger,
) *gin.Engine {
	r := gin.Default()
	// Public
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Lavendera API Gateway is running 🚀",
		})
	})
	api := r.Group("/api/v1")
	// Auth PUBLIC
	auth := api.Group("/auth")
	auth.POST(
		"/register",
		authController.Register,
	)
	auth.POST(
		"/login",
		authController.Login,
	)
	// Protected
	protected := api.Group("")
	protected.Use(
		middleware.JWTMiddleware(jwtConfig),
		middleware.RLSMiddleware(db),
	)
	// protected.GET("/orders", ...)
	// protected.POST("/orders", ...)
	_ = logger
	return r
}
