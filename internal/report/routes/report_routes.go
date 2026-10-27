package routes

import (
	"github.com/Wenell09/lavendera-api/internal/report/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(staffRouter *gin.RouterGroup, controller controller.ReportController) {
	staffRouter.GET("/dashboard", controller.GetDashboard)
	staffRouter.GET("/reports/financial", controller.GetFinancialReport)
}
