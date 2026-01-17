package route

import (
	"github.com/affandisy/goshop/internal/handler"
	"github.com/affandisy/goshop/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler) {
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

		// Category routes (public untuk read, admin untuk write)
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)

			// Admin only
			categories.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		// Product routes (public untuk read, admin untuk write)
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)

			// Admin only
			adminProducts := products.Group("")
			adminProducts.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				adminProducts.POST("", productHandler.Create)
				adminProducts.PUT("/:id", productHandler.Update)
				adminProducts.DELETE("/:id", productHandler.Delete)
				adminProducts.PATCH("/:id/stock", productHandler.UpdateStock)
			}
		}

		// Order routes (protected - need auth)
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetMyOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.POST("/:id/cancel", orderHandler.CancelOrder)

			// Admin only
			orders.GET("/all", middleware.AdminMiddleware(), orderHandler.GetAllOrders)
			orders.PATCH("/:id/status", middleware.AdminMiddleware(), orderHandler.UpdateOrderStatus)
		}
	}
}
