package middleware

import (
	"net/http"
	"time"

	moduleLib "margin-delver/lib"

	"github.com/gin-gonic/gin"
)

func RequestLogger(log *moduleLib.BaseLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.SugarLog().Infof(
			"%s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}

func Recovery(log *moduleLib.BaseLog) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.SugarLog().Errorf("panic recovered: %v", recovered)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Internal Server Error",
			"result":  nil,
		})
	})
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
