package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Для проверки, что сервер живой.
func GetPingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// А как же без этого?
func GetTeapotHandler(c *gin.Context) {
	c.JSON(http.StatusTeapot, gin.H{
		"message": "teapot mode",
	})
}
