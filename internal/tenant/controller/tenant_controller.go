package controller

import "github.com/gin-gonic/gin"

type TenantController interface {
	FindByID(c *gin.Context)
	Update(c *gin.Context)
}
