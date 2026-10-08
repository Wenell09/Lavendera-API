package routes

import (
	"github.com/Wenell09/lavendera-api/internal/service/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(adminRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.ServiceController) {
	adminServices := adminRouter.Group("/services")
	adminServices.POST("", controller.Create)
	adminServices.PATCH("/:id", controller.Update)
	adminServices.DELETE("/:id", controller.Delete)

	staffServices := staffRouter.Group("/services")
	staffServices.GET("", controller.FindAll)
	staffServices.GET("/:id", controller.FindByID)
}
