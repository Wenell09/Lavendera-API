package routes

import (
	"github.com/Wenell09/lavendera-api/internal/service_category/controller"
	"github.com/gin-gonic/gin"
)

func RegisterServiceCategoryRoutes(
	router *gin.RouterGroup,
	controller controller.ServiceCategoryController,
) {
	serviceCategory := router.Group("/service-categories")
	serviceCategory.GET("", controller.FindAll)
	serviceCategory.GET("/:id", controller.FindByID)
	serviceCategory.POST("", controller.Create)
	serviceCategory.PATCH("/:id", controller.Update)
	serviceCategory.DELETE("/:id", controller.Delete)
}
