package routes

import (
	"github.com/Wenell09/lavendera-api/internal/tenant/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.TenantController) {
	tenants := router.Group("/tenants")
	tenants.GET("/:id", controller.FindByID)
	tenants.PATCH("/:id", controller.Update)
}
