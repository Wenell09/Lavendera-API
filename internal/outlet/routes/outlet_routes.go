package routes

import (
	"github.com/Wenell09/lavendera-api/internal/outlet/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.OutletController) {
	outlets := router.Group("/outlets")
	outlets.POST("", controller.Create)
	outlets.GET("", controller.FindAll)
	outlets.GET("/:id", controller.FindByID)
	outlets.PATCH("/:id", controller.Update)
	outlets.DELETE("/:id", controller.Delete)
}
