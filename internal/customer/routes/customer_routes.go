package routes

import (
	"github.com/Wenell09/lavendera-api/internal/customer/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(staffRouter *gin.RouterGroup, controller controller.CustomerController) {
	staffCustomers := staffRouter.Group("/customers")
	staffCustomers.GET("", controller.FindAll)
	staffCustomers.GET("/:id", controller.FindByID)
	staffCustomers.POST("", controller.Create)
	staffCustomers.PATCH("/:id", controller.Update)
	staffCustomers.DELETE("/:id", controller.Delete)
}
