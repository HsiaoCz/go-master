package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	var (
		port   = os.Getenv("PORT")
		router = http.NewServeMux()
	)
	// listen and serve
	logrus.WithFields(logrus.Fields{
		"listen address": port,
	}).Info("the http server is running....")
	
	http.ListenAndServe(port, router)
}
