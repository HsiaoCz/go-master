package main

import (
	"log"
	"net/http"
	"os"

	"github.com/HsiaoCz/go-master/recommend/db"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// init env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// set log
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// init db
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	// init router
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
