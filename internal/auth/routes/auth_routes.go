package routes

import (
	"github.com/Wenell09/lavendera-api/internal/auth/controller"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	controller controller.AuthController,
) {
	auth := router.Group("/auth")
	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
}
