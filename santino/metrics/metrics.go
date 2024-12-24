package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Metrics is a struct that holds all the metrics
	// user metrics
	ActiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chat_active_users",
		Help: "The number of active users currently connected",
	})

	TotalUsers = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_total_users",
		Help: "The total number of registered users",
	})

	// message metrics
	MessagesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_messages_sent_total",
		Help: "The total number of messages sent",
	})

	MessagesPerRoom = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "chat_messages_per_room_total",
		Help: "The total number of messages sent per room",
	}, []string{"room_id"})

	// websocket metrics
	WebSocketConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chat_websocket_connections",
		Help: "The number of active WebSocket connections",
	})

	WebSocketErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_websocket_errors_total",
		Help: "The total number of WebSocket errors",
	})

	// api metrics
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chat_http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method", "status"})

	// database metrics
	DatabaseQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chat_db_query_duration_seconds",
		Help:    "Duration of database queries in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"query_type"})

	// status metrics
	StatusesCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_statuses_created_total",
		Help: "The total number of status updates created",
	})

	// friend metrics
	FriendsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_friends_added_total",
		Help: "The total number of friends added",
	})

	FriendRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_friend_requests_total",
		Help: "The total number of friend requests",
	})

	// private message metrics
	MessagesSentPrivate = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_messages_sent_private_total",
		Help: "The total number of private messages sent",
	})

	MessagesRead = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_messages_read_total",
		Help: "The total number of messages read",
	})

	// room metrics
	RoomsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_rooms_created_total",
		Help: "The total number of rooms created",
	})

	RoomsJoined = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_rooms_joined_total",
		Help: "The total number of rooms joined",
	})

	// avatar metrics
	AvatarsUploaded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_avatars_uploaded_total",
		Help: "The total number of avatars uploaded",
	})

	// status view metrics
	StatusViews = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_status_views_total",
		Help: "The total number of status views",
	})
)
