package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HsiaoCz/go-master/bunt/handlers"
	"github.com/HsiaoCz/go-master/user_test/db"
	logger "github.com/HsiaoCz/go-master/user_test/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.L.WithError(err).Error("Error loading .env file")
		os.Exit(1)
	}

	// Init the database
	if err := db.Init(); err != nil {
		logger.L.Errorf("Error initializing the database: %v", err)
		os.Exit(1)
	}

	var (
		r    = chi.NewRouter()
		port = os.Getenv("PORT")
	)

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Post("/user/signup", handlers.HandleUserSignup)
	r.Post("/user/login", handlers.HandleUserLogin)
	r.Get("/user/{id}", handlers.HandleUserGet)
	r.Put("/user/{id}", handlers.HandleUserUpdate)
	r.Delete("/user/{id}", handlers.HandleUserDelete)

	// post router
	r.Post("/post/create", handlers.HandleCreatePost)
	r.Get("/post/{post_id}", handlers.HandleGetPostByID)
	r.Delete("/post/{post_id}", handlers.HandleDeletePostByID)
	// Start the server

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.L.Errorf("Error starting the server: %v", err)
			os.Exit(1)
		}
	}()

	logger.L.Infof("Server started on port %s", port)

	// Graceful shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.L.Info("Shutting down the server...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Errorf("Error shutting down the server: %v", err)
		os.Exit(1)
	}
}
