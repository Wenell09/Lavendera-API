package routes

import (
	"github.com/Wenell09/lavendera-api/internal/profile/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.ProfileController) {
	profile := router.Group("/profile")
	profile.GET("", controller.GetProfile)
}
