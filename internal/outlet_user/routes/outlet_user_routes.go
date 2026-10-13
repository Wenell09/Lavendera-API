package routes

import (
	"github.com/Wenell09/lavendera-api/internal/outlet_user/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(adminRouter *gin.RouterGroup, staffRouter *gin.RouterGroup, controller controller.OutletUserController) {
	adminOutletUsers := adminRouter.Group("/outlets/:id/users")
	adminOutletUsers.GET("", controller.FindByOutletID)
	adminOutletUsers.POST("", controller.Assign)
	adminOutletUsers.DELETE("/:user_id", controller.Unassign)
	staffOutletUsers := staffRouter.Group("/users/:id/outlets")
	staffOutletUsers.GET("", controller.FindByUserID)
}
