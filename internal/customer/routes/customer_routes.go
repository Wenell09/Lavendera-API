package routes

import (
	"github.com/Wenell09/lavendera-api/internal/customer/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.CustomerController) {
	customers := router.Group("/customers")
	customers.POST("", controller.Create)
	customers.GET("", controller.FindAll)
	customers.GET("/:id", controller.FindByID)
	customers.PATCH("/:id", controller.Update)
	customers.DELETE("/:id", controller.Delete)
}
