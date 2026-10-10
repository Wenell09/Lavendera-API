package controller

import "github.com/gin-gonic/gin"

type OrderController interface {
	CreateCashierOrder(ctx *gin.Context)
	CreatePublicOrder(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	TrackOrderByToken(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
	ConfirmPayment(ctx *gin.Context)
}
