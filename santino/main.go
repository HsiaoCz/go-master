package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HsiaoCz/go-master/santino/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/user/register", handlers.HandleUserRegister).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         os.Getenv("PORT"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error("server error")
			os.Exit(1)
		}
	}()

	logrus.WithFields(logrus.Fields{"port": os.Getenv("PORT")}).Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logrus.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("server shutdown error")
		os.Exit(1)
	}
	logrus.Info("server shutdown")
}
