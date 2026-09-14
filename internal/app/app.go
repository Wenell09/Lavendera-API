package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
	Logger *logrus.Logger
}

func NewApp(
	router *gin.Engine,
	db *gorm.DB,
	logger *logrus.Logger,
) *App {
	return &App{
		Router: router,
		DB:     db,
		Logger: logger,
	}
}
