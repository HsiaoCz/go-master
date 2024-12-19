package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type StreamService struct {
	activeStreams sync.Map // map[string]*Stream
	viewers      sync.Map // map[string]map[*websocket.Conn]bool
	mu           sync.RWMutex
}

type Stream struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Title    string
	Viewers  map[*websocket.Conn]bool
	Done     chan struct{}
	StartedAt time.Time
}

func NewStreamService() *StreamService {
	return &StreamService{}
}

func (s *StreamService) CreateStream(ctx context.Context, userID uuid.UUID, title string) (*Stream, error) {
	stream := &Stream{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Viewers:   make(map[*websocket.Conn]bool),
		Done:      make(chan struct{}),
		StartedAt: time.Now(),
	}

	s.activeStreams.Store(stream.ID.String(), stream)
	return stream, nil
}

func (s *StreamService) EndStream(streamID uuid.UUID) error {
	value, ok := s.activeStreams.LoadAndDelete(streamID.String())
	if !ok {
		return errors.New("stream not found")
	}

	stream := value.(*Stream)
	close(stream.Done)

	// Disconnect all viewers
	s.mu.Lock()
	for conn := range stream.Viewers {
		conn.Close()
		delete(stream.Viewers, conn)
	}
	s.mu.Unlock()

	return nil
}

func (s *StreamService) AddViewer(streamID uuid.UUID, conn *websocket.Conn) error {
	value, ok := s.activeStreams.Load(streamID.String())
	if !ok {
		return errors.New("stream not found")
	}

	stream := value.(*Stream)
	s.mu.Lock()
	stream.Viewers[conn] = true
	s.mu.Unlock()

	return nil
}

func (s *StreamService) RemoveViewer(streamID uuid.UUID, conn *websocket.Conn) {
	value, ok := s.activeStreams.Load(streamID.String())
	if !ok {
		return
	}

	stream := value.(*Stream)
	s.mu.Lock()
	delete(stream.Viewers, conn)
	s.mu.Unlock()
}

func (s *StreamService) GetViewerCount(streamID uuid.UUID) int {
	value, ok := s.activeStreams.Load(streamID.String())
	if !ok {
		return 0
	}

	stream := value.(*Stream)
	s.mu.RLock()
	count := len(stream.Viewers)
	s.mu.RUnlock()

	return count
}
