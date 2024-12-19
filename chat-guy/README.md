# Go Chat Application

A feature-rich chat application built with Go, WebSocket, and PostgreSQL.

## Features

- User authentication (register/login)
- Real-time messaging using WebSocket
- Multiple chat rooms
- Public and private rooms
- Room management (create, join)
- Message persistence
- Modern UI with Tailwind CSS

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Node.js (for development)

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd chat-guy
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up the database:
```bash
# Create a new PostgreSQL database
createdb chatapp

# Apply the schema
psql -d chatapp -f configs/schema.sql
```

4. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials and JWT secret
```

5. Run the application:
```bash
go run cmd/server/main.go
```

The application will be available at `http://localhost:8080`

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── configs/
│   └── schema.sql           # Database schema
├── internal/
│   ├── database/            # Database connection and queries
│   ├── handlers/            # HTTP and WebSocket handlers
│   ├── middleware/          # Authentication middleware
│   └── models/              # Data models
├── static/                  # Frontend assets
│   └── index.html          # Main frontend file
├── .env.example            # Environment variables template
├── go.mod                  # Go module file
└── README.md               # Project documentation
```

## API Endpoints

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/rooms` - List all accessible rooms
- `POST /api/rooms` - Create a new room
- `POST /api/rooms/{id}/join` - Join a room
- `WS /api/ws` - WebSocket endpoint for real-time chat

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License.
