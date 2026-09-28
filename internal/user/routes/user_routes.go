package routes

import (
	"github.com/Wenell09/lavendera-api/internal/user/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.UserController) {
	users := router.Group("/users")
	users.POST("", controller.Create)
	users.GET("", controller.FindAll)
	users.GET("/:id", controller.FindByID)
	users.PATCH("/:id", controller.Update)
	users.DELETE("/:id", controller.Delete)
}
