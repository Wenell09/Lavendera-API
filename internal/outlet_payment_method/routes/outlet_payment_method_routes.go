package routes

import (
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.OutletPaymentMethodController) {
	outletPaymentMethods := router.Group("/outlets/:id/payment-methods")
	outletPaymentMethods.POST("", controller.Create)
	outletPaymentMethods.GET("", controller.FindAll)
	outletPaymentMethods.GET("/:payment_method_id", controller.FindByID)
	outletPaymentMethods.PATCH("/:payment_method_id", controller.Update)
	outletPaymentMethods.DELETE("/:payment_method_id", controller.Delete)
}
