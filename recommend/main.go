package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/HsiaoCz/go-master/recommend/db"
	"github.com/HsiaoCz/go-master/recommend/handlers"
	"github.com/HsiaoCz/go-master/recommend/handlers/middlewares"
	"github.com/HsiaoCz/go-master/recommend/mod"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	// set log
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	// init env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// init db
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	// init redis

	go func() {
		count, err := strconv.Atoi(os.Getenv("DBCOUNT"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error message": err,
			}).Error("dbcount must be int")
			os.Exit(1)
		}
		db.InitRedis(os.Getenv("REDISURL"), os.Getenv("PASSWD"), count)
	}()

	// connect mongo db
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
		port         = os.Getenv("PORT")
		router       = http.NewServeMux()
		userData     = mod.UserModInit(db.Get())
		userHandlers = handlers.UserHandlersInit(userData)
	)

	{
		// user handlefunc
		router.HandleFunc("POST /api/v1/user", handlers.TransferHandlerfunc(userHandlers.HandleCreateUser))
		router.HandleFunc("GET /api/v1/user/{user_id}", handlers.TransferHandlerfunc(userHandlers.HandleCreateUser))
		router.HandleFunc("DELETE /api/v1/user", middlewares.JwtMiddleware(handlers.TransferHandlerfunc(userHandlers.HandleDeleteUserByID)))

		// record
	}

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
