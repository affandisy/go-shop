package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := recover(); err != nil {
			log.Printf("PANIC: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal server error",
				"error":   "Unexpected Error",
			})

			c.Abort()
		}
		c.Next()
	}
}
