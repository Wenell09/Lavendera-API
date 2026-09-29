package app

import (
	"net/http"

	authController "github.com/Wenell09/lavendera-api/internal/auth/controller"
	authRoutes "github.com/Wenell09/lavendera-api/internal/auth/routes"
	customerController "github.com/Wenell09/lavendera-api/internal/customer/controller"
	customerRoutes "github.com/Wenell09/lavendera-api/internal/customer/routes"
	discountController "github.com/Wenell09/lavendera-api/internal/discount/controller"
	discountRoutes "github.com/Wenell09/lavendera-api/internal/discount/routes"
	OutletController "github.com/Wenell09/lavendera-api/internal/outlet/controller"
	OutletRoutes "github.com/Wenell09/lavendera-api/internal/outlet/routes"
	outletUserController "github.com/Wenell09/lavendera-api/internal/outlet_user/controller"
	outletUserRoutes "github.com/Wenell09/lavendera-api/internal/outlet_user/routes"
	serviceController "github.com/Wenell09/lavendera-api/internal/service/controller"
	serviceRoutes "github.com/Wenell09/lavendera-api/internal/service/routes"
	serviceCategoryController "github.com/Wenell09/lavendera-api/internal/service_category/controller"
	serviceCategoryRoutes "github.com/Wenell09/lavendera-api/internal/service_category/routes"
	"github.com/Wenell09/lavendera-api/internal/shared/config"
	"github.com/Wenell09/lavendera-api/internal/shared/middleware"
	userController "github.com/Wenell09/lavendera-api/internal/user/controller"
	userRoutes "github.com/Wenell09/lavendera-api/internal/user/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(
	authController authController.AuthController,
	serviceCategoryController serviceCategoryController.ServiceCategoryController,
	serviceController serviceController.ServiceController,
	discountController discountController.DiscountController,
	customerController customerController.CustomerController,
	outletController OutletController.OutletController,
	userController userController.UserController,
	outletUserController outletUserController.OutletUserController,
	jwtConfig config.JWTConfig,
	logger *logrus.Logger,
) *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Lavendera API Gateway is running 🚀",
		})
	})
	api := r.Group("/api/v1")
	// Route public
	publicGroup := api.Group("")
	authRoutes.RegisterRoutes(publicGroup, authController)
	// Route protected
	protectedGroup := api.Group("")
	protectedGroup.Use(middleware.JWTMiddleware(jwtConfig))
	// Registrasi fitur yang membutuhkan autentikasi
	// Route khusus admin
	adminGroup := protectedGroup.Group("")
	adminGroup.Use(middleware.RequireRole("ADMIN"))
	OutletRoutes.RegisterRoutes(adminGroup, outletController)
	serviceCategoryRoutes.RegisterRoutes(adminGroup, serviceCategoryController)
	serviceRoutes.RegisterRoutes(adminGroup, serviceController)
	discountRoutes.RegisterRoutes(adminGroup, discountController)
	customerRoutes.RegisterRoutes(adminGroup, customerController)
	userRoutes.RegisterRoutes(adminGroup, userController)
	outletUserRoutes.RegisterRoutes(adminGroup, outletUserController)
	// route khusus staff
	return r
}
