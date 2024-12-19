package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// User metrics
	ActiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chat_active_users",
		Help: "The number of active users currently connected",
	})

	TotalUsers = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_total_users",
		Help: "The total number of registered users",
	})

	// Message metrics
	MessagesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_messages_sent_total",
		Help: "The total number of messages sent",
	})

	MessagesPerRoom = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "chat_messages_per_room_total",
		Help: "The total number of messages sent per room",
	}, []string{"room_id"})

	// WebSocket metrics
	WebSocketConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chat_websocket_connections",
		Help: "The number of active WebSocket connections",
	})

	WebSocketErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_websocket_errors_total",
		Help: "The total number of WebSocket errors",
	})

	// API metrics
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chat_http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method", "status"})

	// Database metrics
	DatabaseQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chat_db_query_duration_seconds",
		Help:    "Duration of database queries in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"query_type"})

	// Status metrics
	StatusesCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_statuses_created_total",
		Help: "The total number of status updates created",
	})

	StatusViews = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chat_status_views_total",
		Help: "The total number of status views",
	})
)
