package routes

import (
	"github.com/Wenell09/lavendera-api/internal/service/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.ServiceController) {
	services := router.Group("/services")
	services.POST("", controller.Create)
	services.GET("", controller.FindAll)
	services.GET("/:id", controller.FindByID)
	services.PATCH("/:id", controller.Update)
	services.DELETE("/:id", controller.Delete)
}
