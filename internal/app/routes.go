package app

import (
	"net/http"

	authController "github.com/Wenell09/lavendera-api/internal/auth/controller"
	authRoutes "github.com/Wenell09/lavendera-api/internal/auth/routes"
	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/Wenell09/lavendera-api/internal/middleware"
	serviceCategoryController "github.com/Wenell09/lavendera-api/internal/service_category/controller"
	serviceCategoryAuth "github.com/Wenell09/lavendera-api/internal/service_category/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewRouter(
	authController authController.AuthController,
	serviceCategoryController serviceCategoryController.ServiceCategoryController,
	jwtConfig config.JWTConfig,
	db *gorm.DB,
	logger *logrus.Logger,
) *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Lavendera API Gateway is running 🚀",
		})
	})
	api := r.Group("/api/v1")
	// Route public
	publicGroup := api.Group("")
	authRoutes.RegisterAuthRoutes(publicGroup, authController)
	// Route protected
	protectedGroup := api.Group("")
	protectedGroup.Use(
		middleware.JWTMiddleware(jwtConfig),
		middleware.RLSMiddleware(db),
	)
	// Registrasi fitur yang membutuhkan autentikasi
	serviceCategoryAuth.RegisterServiceCategoryRoutes(protectedGroup, serviceCategoryController)
	_ = logger
	return r
}
