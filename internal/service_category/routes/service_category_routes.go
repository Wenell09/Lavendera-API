package routes

import (
	"github.com/Wenell09/lavendera-api/internal/service_category/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(adminRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.ServiceCategoryController) {
	adminServiceCategory := adminRouter.Group("/service-categories")
	adminServiceCategory.POST("", controller.Create)
	adminServiceCategory.PATCH("/:id", controller.Update)
	adminServiceCategory.DELETE("/:id", controller.Delete)

	staffServiceCategory := staffRouter.Group("/service-categories")
	staffServiceCategory.GET("", controller.FindAll)
	staffServiceCategory.GET("/:id", controller.FindByID)

}
