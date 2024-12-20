package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shopping-mall/internal/handlers"
	"github.com/shopping-mall/internal/repository/mongodb"
	"github.com/shopping-mall/internal/service"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	// Initialize MongoDB repository
	if err := godotenv.Load("../../.env"); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error loading env file")
		os.Exit(1)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	repo, err := mongodb.NewMongoRepository(mongoURI)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to connect to MongoDB")
		os.Exit(1)
	}

	// Initialize service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	svc := service.NewService(repo, repo, repo, jwtSecret)

	// Initialize handlers
	h := handlers.NewHandler(svc)

	// Initialize router
	router := gin.Default()

	// Public routes
	router.POST("/api/register", h.Register)
	router.POST("/api/login", h.Login)
	router.GET("/api/products", h.ListProducts)
	router.GET("/api/products/:id", h.GetProduct)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(handlers.AuthMiddleware(jwtSecret))
	{
		// User routes
		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders", h.ListUserOrders)
		protected.GET("/orders/:id", h.GetOrder)

		// Admin routes
		admin := protected.Group("")
		admin.Use(handlers.AdminMiddleware())
		{
			admin.POST("/products", h.CreateProduct)
			admin.PUT("/orders/:id/status", h.UpdateOrderStatus)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to start server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Server forced to shutdown")
		os.Exit(1)
	}

	logrus.Info("Server exiting")
}
