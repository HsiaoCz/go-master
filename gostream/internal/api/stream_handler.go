package api

import (
	"gostream/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type StreamHandler struct {
	streamService *service.StreamService
	upgrader      websocket.Upgrader
}

func NewStreamHandler(streamService *service.StreamService) *StreamHandler {
	return &StreamHandler{
		streamService: streamService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, implement proper origin checking
			},
		},
	}
}

func (h *StreamHandler) StartStream(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, get userID from authenticated session
	userID := uuid.New()

	stream, err := h.streamService.CreateStream(c.Request.Context(), userID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stream_id": stream.ID,
		"title":     stream.Title,
	})
}

func (h *StreamHandler) EndStream(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	if err := h.streamService.EndStream(streamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stream ended"})
}

func (h *StreamHandler) WatchStream(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if err := h.streamService.AddViewer(streamID, conn); err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer h.streamService.RemoveViewer(streamID, conn)

	// Handle WebSocket connection
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.CloseMessage {
			break
		}
	}
}

func (h *StreamHandler) GetStreamInfo(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	viewerCount := h.streamService.GetViewerCount(streamID)
	c.JSON(http.StatusOK, gin.H{
		"stream_id":    streamID,
		"viewer_count": viewerCount,
	})
}
