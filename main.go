package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var port string = ":8000"

type body struct {
	Detail string `json:"detail"`
}

var (
	interval = 2 // インターバル時間
	burst    = 2 // 同時接続数
	limit    = rate.NewLimiter(rate.Every(time.Duration(interval)*time.Second), burst)
)

func example() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, body{
			Detail: http.StatusText(http.StatusOK),
		})
	}
}

func rateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit.Allow() == false {
			c.JSON(http.StatusTooManyRequests, body{
				Detail: http.StatusText(http.StatusTooManyRequests),
			})
			c.Abort()
		}
	}
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.Use(rateLimiter())

	router.GET("/example", example())

	return router
}

func main() {
	router := initRouter()
	router.Run(port)
}
