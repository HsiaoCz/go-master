package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HsiaoCz/go-master/bunt/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Errorf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	// Init the database
	if err := db.Init(); err != nil {
		logrus.Errorf("Error initializing the database: %v", err)
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

	// Start the server

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("Error starting the server: %v", err)
			os.Exit(1)
		}
	}()

	logrus.Infof("Server started on port %s", port)

	// Graceful shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Info("Shutting down the server...")

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Error shutting down the server: %v", err)
		os.Exit(1)
	}

}
