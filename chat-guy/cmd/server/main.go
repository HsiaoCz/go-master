package main

import (
	"log"
	"net/http"
	"os"

	"chat-guy/internal/database"
	"chat-guy/internal/handlers"
	"chat-guy/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// Add metrics middleware
	r.Use(middleware.MetricsMiddleware)

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Public routes
	r.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// User profile routes
	api.HandleFunc("/user/profile", handlers.GetUserProfile).Methods("GET")

	// Avatar routes
	api.HandleFunc("/avatar", handlers.UploadAvatar).Methods("POST")
	r.PathPrefix("/uploads/avatars/").Handler(http.StripPrefix("/uploads/avatars/", http.HandlerFunc(handlers.ServeAvatar)))

	// Room routes
	api.HandleFunc("/rooms", handlers.CreateRoom).Methods("POST")
	api.HandleFunc("/rooms", handlers.GetRooms).Methods("GET")
	api.HandleFunc("/rooms/{id}/join", handlers.JoinRoomHTTP).Methods("POST")

	// Friend routes
	api.HandleFunc("/friends", handlers.GetFriends).Methods("GET")
	api.HandleFunc("/friends/requests", handlers.GetFriendRequests).Methods("GET")
	api.HandleFunc("/friends/{username}", handlers.AddFriend).Methods("POST")
	api.HandleFunc("/friends/{id}/accept", handlers.AcceptFriendRequest).Methods("POST")

	// Private message routes
	api.HandleFunc("/messages/{id}", handlers.GetPrivateMessages).Methods("GET")
	api.HandleFunc("/messages/{id}/send", handlers.SendPrivateMessage).Methods("POST")

	// Status routes
	api.HandleFunc("/status", handlers.CreateStatus).Methods("POST")
	api.HandleFunc("/status/friends", handlers.GetFriendsStatuses).Methods("GET")
	api.HandleFunc("/status/{id}/view", handlers.ViewStatus).Methods("POST")
	r.PathPrefix("/uploads/status/").Handler(http.StripPrefix("/uploads/status/", http.HandlerFunc(handlers.ServeAvatar)))

	// WebSocket route
	api.HandleFunc("/ws", handlers.HandleWebSocket)

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s\n", port)
	log.Printf("Metrics available on http://localhost:%s/metrics\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
