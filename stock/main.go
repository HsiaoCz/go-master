package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGOURL")))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := client.Ping(ctx, &readpref.ReadPref{}); err != nil {
			log.Fatal(err)
		}
	}()

	var (
		port   = os.Getenv("PORT")
		router = http.NewServeMux()
	)
	
	server := http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  time.Millisecond * 1500,
		WriteTimeout: time.Millisecond * 1500,
	}

	go func() {
		logrus.WithFields(logrus.Fields{
			"listen address": port,
		}).Info("the http server is running....")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	log.Println("shutting down the server....")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server gracefully shut down....")
}
