package routes

import (
	"github.com/Wenell09/lavendera-api/internal/order/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(publicRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.OrderController) {
	// Public routes (no JWT required)
	publicRouter.POST("/public/orders", controller.CreatePublicOrder)
	publicRouter.GET("/public/orders/track/:token", controller.TrackOrderByToken)

	// Staff routes (Staff & Admin)
	staffOrders := staffRouter.Group("/orders")
	staffOrders.POST("", controller.CreateCashierOrder)
	staffOrders.GET("", controller.FindAll)
	staffOrders.GET("/:id", controller.FindByID)
	staffOrders.PATCH("/:id/status", controller.UpdateStatus)
	staffOrders.PATCH("/:id/cancel", controller.CancelOrder)
	staffOrders.PATCH("/:id/payments/:payment_id/confirm", controller.ConfirmPayment)
	staffOrders.PATCH("/:id/cash-payment", controller.InputCashPayment)

}
