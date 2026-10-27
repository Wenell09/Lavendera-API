package controller

import "github.com/gin-gonic/gin"

type ReportController interface {
	GetDashboard(ctx *gin.Context)
	GetFinancialReport(ctx *gin.Context)
}
