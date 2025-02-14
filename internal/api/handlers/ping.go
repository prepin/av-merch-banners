package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

// Для тестирования отсечки слишком длинных запросов по таймауту.
func GetSleepHandler(c *gin.Context) {
	sleepParam := c.DefaultQuery("timeout", "100")
	sleep, err := strconv.Atoi(sleepParam)
	if err != nil {
		sleep = 100
	}

	select {
	case <-time.After(time.Duration(sleep) * time.Millisecond):
		if c.Request.Context().Err() != nil {
			return
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf("slept %d ms", sleep)})

	case <-c.Request.Context().Done():
		return
	}
}
