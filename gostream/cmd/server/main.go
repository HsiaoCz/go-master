package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gostream/configs"
	"gostream/internal/api"
	"gostream/internal/middleware"
	"gostream/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()
	
	// Initialize services
	streamService := service.NewStreamService()
	
	// Initialize handlers
	streamHandler := api.NewStreamHandler(streamService)
	
	// Setup router
	router := gin.Default()
	
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	
	// Setup routes
	setupRoutes(router, streamHandler)

	srv := &http.Server{
		Addr:    ":" + config.ServerPort,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRoutes(r *gin.Engine, streamHandler *api.StreamHandler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public routes
		api.GET("/streams", streamHandler.GetStreamInfo)
		
		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/streams", streamHandler.StartStream)
			protected.DELETE("/streams/:id", streamHandler.EndStream)
			protected.GET("/streams/:id/watch", streamHandler.WatchStream)
		}
	}
}
