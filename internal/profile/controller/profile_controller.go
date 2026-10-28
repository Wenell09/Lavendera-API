package controller

import "github.com/gin-gonic/gin"

type ProfileController interface {
	GetProfile(c *gin.Context)
}
