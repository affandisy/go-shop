package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/affandisy/goshop/cmd/route"
	"github.com/affandisy/goshop/internal/handler"
	"github.com/affandisy/goshop/internal/middleware"
	"github.com/affandisy/goshop/internal/repository"
	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/cache"
	"github.com/affandisy/goshop/pkg/config"
	"github.com/affandisy/goshop/pkg/database"
	"github.com/affandisy/goshop/pkg/payment"
	"github.com/affandisy/goshop/pkg/redis"
	"github.com/affandisy/goshop/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Main File")

	cfg := config.Load("pkg/config/config.yaml")

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init jwt
	utils.InitJWt(cfg.AuthSecret)

	// Connect to Database
	if err := database.Connect(cfg.DatabaseURI, cfg.Environment == "production"); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	// Auto Migration
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	// Seed Data
	if err := database.SeedData(); err != nil {
		log.Fatalf("Seeding data failed: %v", err)
	}

	// Connect to Redis
	if err := redis.Connect(cfg.RedisURI, cfg.RedisPassword, cfg.RedisDB); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer redis.Close()

	redisClient := redis.GetClient()
	cacheService := cache.NewCacheService(redisClient)

	// midtrans client
	midtransClient := payment.NewMidtransClient(cfg.MidtransServerKey, cfg.MidtransClientKey, cfg.MidtransEnvironment)

	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo, cacheService)
	productService := service.NewProductService(productRepo, categoryRepo, cacheService)
	orderService := service.NewOrderService(orderRepo, productRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo, midtransClient)

	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	cacheHandler := handler.NewCacheHandler(cacheService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	route.SetupRoutes(router, userHandler, categoryHandler, productHandler, orderHandler, paymentHandler, cacheHandler)

	log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Health check: http://localhost:%s/health", cfg.HTTPPort)
	log.Printf("API Base URL: http://localhost:%s/api/v1", cfg.HTTPPort)

	// Graceful shutdown
	go func() {
		if err := router.Run(":" + cfg.HTTPPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
