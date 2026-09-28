package controller

import "github.com/gin-gonic/gin"

type OutletUserController interface {
	Assign(c *gin.Context)
	Unassign(c *gin.Context)
	FindByOutletID(c *gin.Context)
}
