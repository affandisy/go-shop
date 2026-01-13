package route

import (
	"github.com/affandisy/goshop/internal/handler"
	"github.com/affandisy/goshop/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"message":  "GoShop API is running",
			"database": "connected",
			"redis":    "connected",
		})
	})

	v1 := router.Group("/api/v1")
	{
		// Ping endpoint
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Auth routes (public - tidak perlu auth)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// User routes (protected - perlu auth)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)

			// Admin only routes
			users.GET("", middleware.AdminMiddleware(), userHandler.GetUsers)
		}
	}
}
