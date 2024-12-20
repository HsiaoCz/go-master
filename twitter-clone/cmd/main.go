package main

import (
	"twitter-clone/internal/config"
	"twitter-clone/internal/database"
	"twitter-clone/internal/handlers"
	"twitter-clone/internal/middleware"
	"twitter-clone/internal/services"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	// Initialize services and handlers
	authService := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(authService)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Routes
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes example
	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware)

	// Initialize tweet handler
	tweetService := services.NewTweetService()
	tweetHandler := handlers.NewTweetHandler(tweetService)

	// Tweet routes
	tweet := protected.Group("/tweet")
	tweet.POST("", tweetHandler.CreateTweet)
	tweet.GET("/:id", tweetHandler.GetTweet)
	tweet.PUT("/:id", tweetHandler.UpdateTweet)
	tweet.DELETE("/:id", tweetHandler.DeleteTweet)
	tweet.GET("/user/:user_id", tweetHandler.GetUserTweets)
	// Timeline routes
	tweet.GET("/timeline", tweetHandler.GetHomeTimeline)
	tweet.GET("/user/:user_id/timeline", tweetHandler.GetUserTimeline)
	tweet.GET("/:id/replies", tweetHandler.GetTweetReplies)
	// Initialize interaction services and handler
	likeService := services.NewLikeService()
	retweetService := services.NewRetweetService()
	interactionHandler := handlers.NewInteractionHandler(likeService, retweetService)

	// Interaction routes
	tweet.POST("/:id/like", interactionHandler.LikeTweet)
	tweet.DELETE("/:id/like", interactionHandler.UnlikeTweet)
	tweet.GET("/:id/likes", interactionHandler.GetTweetLikes)
	tweet.POST("/:id/retweet", interactionHandler.Retweet)
	tweet.DELETE("/:id/retweet", interactionHandler.UndoRetweet)
	tweet.GET("/:id/retweets", interactionHandler.GetTweetRetweets)

	// Initialize follow service and handler
	followService := services.NewFollowService()
	followHandler := handlers.NewFollowHandler(followService)

	// Follow routes
	users := protected.Group("/users")
	users.POST("/:id/follow", followHandler.FollowUser)
	users.DELETE("/:id/follow", followHandler.UnfollowUser)
	users.GET("/:id/followers", followHandler.GetFollowers)
	users.GET("/:id/following", followHandler.GetFollowing)
	users.GET("/:id/follow-status", followHandler.CheckFollowStatus)

	// Initialize profile service and handler
	profileService := services.NewProfileService()
	profileHandler := handlers.NewProfileHandler(profileService)

	// Profile routes
	profile := protected.Group("/profile")
	profile.POST("/background", profileHandler.UpdateBackground)
	profile.PUT("/", profileHandler.UpdateProfile)

	// Serve static files
	e.Static("/uploads", "uploads")
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
