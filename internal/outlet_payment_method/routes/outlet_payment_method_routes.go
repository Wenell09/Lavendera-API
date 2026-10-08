package routes

import (
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(adminRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.OutletPaymentMethodController) {
	adminOutletPaymentMethods := adminRouter.Group("/outlets/:id/payment-methods")
	adminOutletPaymentMethods.POST("", controller.Create)
	adminOutletPaymentMethods.PATCH("/:payment_method_id", controller.Update)
	adminOutletPaymentMethods.DELETE("/:payment_method_id", controller.Delete)

	staffOutletPaymentMethods := staffRouter.Group("/outlets/:id/payment-methods")
	staffOutletPaymentMethods.GET("", controller.FindAll)
	staffOutletPaymentMethods.GET("/:payment_method_id", controller.FindByID)
}
