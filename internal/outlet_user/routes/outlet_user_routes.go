package routes

import (
	"github.com/Wenell09/lavendera-api/internal/outlet_user/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.OutletUserController) {
	outletUsers := router.Group("/outlets/:id/users")
	{
		outletUsers.GET("", controller.FindByOutletID)
		outletUsers.POST("", controller.Assign)
		outletUsers.DELETE("/:user_id", controller.Unassign)
	}
}
