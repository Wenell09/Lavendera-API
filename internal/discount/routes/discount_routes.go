package routes

import (
	"github.com/Wenell09/lavendera-api/internal/discount/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(adminRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.DiscountController) {
	adminDiscounts := adminRouter.Group("/discounts")
	adminDiscounts.POST("", controller.Create)
	adminDiscounts.PATCH("/:id", controller.Update)
	adminDiscounts.DELETE("/:id", controller.Delete)

	staffDiscounts := staffRouter.Group("/discounts")
	staffDiscounts.GET("", controller.FindAll)
	staffDiscounts.GET("/:id", controller.FindByID)
}
