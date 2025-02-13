package main

import (
	"os"

	"github.com/HsiaoCz/go-master/teml/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.L.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to load env file")
		os.Exit(1)
	}

	
}
