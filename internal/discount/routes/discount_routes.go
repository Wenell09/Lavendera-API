package routes

import (
	"github.com/Wenell09/lavendera-api/internal/discount/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.DiscountController) {
	discounts := router.Group("/discounts")
	discounts.POST("", controller.Create)
	discounts.GET("", controller.FindAll)
	discounts.GET("/:id", controller.FindByID)
	discounts.PATCH("/:id", controller.Update)
	discounts.DELETE("/:id", controller.Delete)
}
