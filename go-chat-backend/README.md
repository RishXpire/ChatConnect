# ChatConnect Backend

Real-time chat application backend built with Go, WebSockets, and MongoDB.

## Features

- RESTful API endpoints for authentication and room management
- WebSocket server for real-time messaging
- MongoDB integration for data persistence
- bcrypt password hashing
- Concurrent room management with Go routines

## API Documentation

### Authentication Endpoints

#### Register User
```http
POST /register
Content-Type: application/json

{
  "username": "user123",
  "password_hash": "securepassword"
}
```

**Response:**
```json
{
  "message": "User registered successfully"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "username": "user123",
  "password_hash": "securepassword"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "username": "user123"
}
```

### Room Management

#### List Rooms
```http
GET /rooms
```

**Response:**
```json
[
  {
    "name": "General",
    "creator": "user123",
    "created_at": "2025-11-21T12:00:00Z"
  }
]
```

#### Delete Room
```http
DELETE /rooms?name=General&username=user123
```

**Response:**
```json
{
  "message": "Room deleted successfully"
}
```

### WebSocket

#### Connect to Room
```
WS /ws/{roomName}?username={username}
```

**Send Message:**
```json
{
  "room_id": "General",
  "username": "user123",
  "content": "Hello, world!"
}
```

**Receive Message:**
```json
{
  "room_id": "General",
  "username": "user123",
  "content": "Hello, world!",
  "timestamp": "2025-11-21T12:00:00Z"
}
```

## Running the Server

```bash
# Install dependencies
go mod download

# Run server
go run main.go
```

Server runs on `http://localhost:8000`

## Environment Variables

Create a `.env` file (see `.env.example`):

```
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/
PORT=8000
```

## Building

```bash
# Build binary
go build -o chatserver main.go

# Run binary
./chatserver
```

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Dependencies

- `gorilla/websocket` - WebSocket implementation
- `go.mongodb.org/mongo-driver` - MongoDB driver
- `golang.org/x/crypto` - bcrypt password hashing

## Project Structure

```
go-chat-backend/
├── main.go              # HTTP server and API handlers
├── models/
│   └── models.go        # Data structures
├── websocket/
│   ├── hub.go          # Hub manager
│   └── client.go       # Client handler
├── go.mod
└── go.sum
```

## Security Notes

- Passwords are hashed with bcrypt (cost 10)
- No plaintext passwords stored
- Authorization checks for room deletion
- CORS configured (update for production)

## Production Deployment

1. Set environment variables
2. Build binary: `go build`
3. Deploy to server (Railway, Render, AWS, etc.)
4. Configure reverse proxy (nginx)
5. Enable HTTPS
6. Update CORS settings
