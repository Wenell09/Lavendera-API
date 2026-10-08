package routes

import (
	"github.com/Wenell09/lavendera-api/internal/customer/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(adminRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.CustomerController) {
	adminCustomers := adminRouter.Group("/customers")
	adminCustomers.POST("", controller.Create)
	adminCustomers.PATCH("/:id", controller.Update)
	adminCustomers.DELETE("/:id", controller.Delete)

	staffCustomers := staffRouter.Group("/customers")
	staffCustomers.GET("", controller.FindAll)
	staffCustomers.GET("/:id", controller.FindByID)
}
