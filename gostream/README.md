# GoStream - A Golang Streaming Platform

A real-time streaming platform built with Go, featuring WebSocket-based live streaming capabilities.

## Features

- Live streaming with WebSocket support
- User authentication
- Stream management (create, end, watch)
- Real-time viewer counting
- Graceful shutdown handling

## Prerequisites

- Go 1.21 or higher
- Git

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd gostream
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default.

## API Endpoints

### Public Endpoints
- `GET /health` - Health check
- `GET /api/streams` - Get information about active streams

### Protected Endpoints (Requires Authentication)
- `POST /api/streams` - Start a new stream
- `DELETE /api/streams/:id` - End a stream
- `GET /api/streams/:id/watch` - Watch a stream (WebSocket)

## Authentication

Include a Bearer token in the Authorization header:
```
Authorization: Bearer your-token-here
```

## Environment Variables

- `SERVER_PORT` - Server port (default: 8080)
- `JWT_SECRET` - Secret key for JWT authentication

## License

MIT
