package controller

import (
	"github.com/gin-gonic/gin"
)

type OutletPaymentMethodController interface {
	Create(c *gin.Context)
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
