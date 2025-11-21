# ChatConnect - Real-Time Chat Application

A full-stack real-time chat application built with Go, MongoDB, and React. Features include user authentication, persistent message storage, room management, and WebSocket-based real-time communication.

## Features

- **User Authentication**: Secure registration and login with bcrypt password hashing
- **Real-Time Messaging**: WebSocket-based instant message delivery
- **Persistent Storage**: Messages, rooms, and users stored in MongoDB
- **Room Management**: Create, join, and delete chat rooms
- **Authorization**: Room creators can manage their rooms
- **Modern UI**: Glassmorphism design with responsive layout
- **Message History**: Automatic loading of previous messages when joining rooms

## Tech Stack

### Backend
- **Go 1.19+**: High-performance backend server
- **MongoDB**: NoSQL database for data persistence
- **Gorilla WebSocket**: Real-time bidirectional communication
- **bcrypt**: Secure password hashing

### Frontend
- **React 18**: Modern component-based UI
- **WebSocket API**: Native browser WebSocket support
- **CSS3**: Custom styling with glassmorphism effects

## Project Structure

```
ChatApp/
├── go-chat-backend/          # Go backend server
│   ├── main.go               # HTTP server and API endpoints
│   ├── models/
│   │   └── models.go         # Data models (User, Room, Message)
│   ├── websocket/
│   │   ├── hub.go            # WebSocket hub manager
│   │   └── client.go         # WebSocket client handling
│   ├── go.mod                # Go module dependencies
│   └── go.sum
│
└── chat-frontend-react/      # React frontend application
    ├── public/
    ├── src/
    │   ├── App.js            # Main application component
    │   ├── App.css           # Application styles
    │   ├── components/
    │   │   └── InputForm.js  # Message input component
    │   └── index.js          # Application entry point
    ├── package.json
    └── package-lock.json
```

## Prerequisites

- **Go**: Version 1.19 or higher
- **Node.js**: Version 16 or higher
- **npm**: Version 8 or higher
- **MongoDB**: Local instance or MongoDB Atlas account

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/ChatConnect.git
cd ChatConnect
```

### 2. Backend Setup

```bash
cd go-chat-backend

# Install dependencies
go mod download

# Update MongoDB connection string in main.go (line 86)
# Replace with your MongoDB URI

# Run the server
go run main.go
```

The backend server will start on `http://localhost:8000`

### 3. Frontend Setup

```bash
cd chat-frontend-react

# Install dependencies
npm install

# Start development server
npm start
```

The frontend will open at `http://localhost:3000`

## Configuration

### MongoDB Connection

Update the connection string in `go-chat-backend/main.go`:

```go
connectionString := "your-mongodb-connection-string"
```

For MongoDB Atlas:
```
mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
```

For local MongoDB:
```
mongodb://localhost:27017
```

### CORS Configuration

The backend is configured to accept requests from any origin (`*`). For production, update CORS settings in `main.go`:

```go
w.Header().Set("Access-Control-Allow-Origin", "https://yourdomain.com")
```

## API Endpoints

### Authentication

- `POST /register` - Register a new user
  - Body: `{ "username": "string", "password_hash": "string" }`
  - Returns: `{ "message": "User registered successfully" }`

- `POST /login` - Authenticate user
  - Body: `{ "username": "string", "password_hash": "string" }`
  - Returns: `{ "message": "Login successful", "username": "string" }`

### Rooms

- `GET /rooms` - List all rooms
  - Returns: `[{ "name": "string", "creator": "string", "created_at": "timestamp" }]`

- `DELETE /rooms?name={roomName}&username={username}` - Delete a room (creator only)
  - Returns: `{ "message": "Room deleted successfully" }`

### WebSocket

- `WS /ws/{roomName}?username={username}` - Connect to room
  - Send: `{ "room_id": "string", "username": "string", "content": "string" }`
  - Receive: `{ "room_id": "string", "username": "string", "content": "string", "timestamp": "ISO8601" }`

## Usage

### Registration and Login

1. Navigate to `http://localhost:3000`
2. Click "New here? Create Account"
3. Enter username and password
4. Click "Sign Up"
5. Login with your credentials

### Creating and Joining Rooms

1. After login, you'll see the Lobby
2. Enter a room name and click "Create & Join"
3. Or select an existing room and click "Join"

### Sending Messages

1. Type your message in the input field
2. Press Enter or click Send
3. Messages appear instantly for all connected users

### Deleting Rooms

1. In the Lobby, find rooms you created
2. Click the delete icon (trash can) next to your room
3. Confirm deletion
4. All users will be disconnected from the room

## Database Schema

### Users Collection

```javascript
{
  "_id": ObjectId,
  "username": String,
  "password_hash": String (bcrypt hash)
}
```

### Rooms Collection

```javascript
{
  "_id": ObjectId,
  "name": String,
  "creator": String,
  "created_at": ISODate
}
```

### Messages Collection

```javascript
{
  "_id": ObjectId,
  "room_id": String,
  "username": String,
  "content": String,
  "timestamp": ISODate
}
```

## Security Considerations

- Passwords are hashed using bcrypt before storage
- No plaintext passwords are stored in the database
- Authorization checks ensure only room creators can delete rooms
- CORS should be configured for production environments
- Consider implementing JWT tokens for session management
- Add rate limiting for production deployment

## Development

### Running Tests

Backend:
```bash
cd go-chat-backend
go test ./...
```

Frontend:
```bash
cd chat-frontend-react
npm test
```

### Building for Production

Backend:
```bash
cd go-chat-backend
go build -o chatserver main.go
./chatserver
```

Frontend:
```bash
cd chat-frontend-react
npm run build
# Serve the build directory with a static file server
```

## Deployment

### Backend Deployment

Recommended platforms:
- **Railway**: Easy Go deployment
- **Render**: Free tier available
- **Heroku**: Container deployment
- **AWS EC2**: Full control

### Frontend Deployment

Recommended platforms:
- **Vercel**: Optimized for React
- **Netlify**: Easy static hosting
- **GitHub Pages**: Free hosting
- **AWS S3 + CloudFront**: Scalable solution

### Environment Variables

For production, use environment variables:

```bash
export MONGO_URI="your-mongodb-connection-string"
export PORT="8000"
export JWT_SECRET="your-secret-key"
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -m 'Add YourFeature'`)
4. Push to the branch (`git push origin feature/YourFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Rishabh**
- GitHub: [@yourusername](https://github.com/yourusername)
- LinkedIn: [Your Name](https://linkedin.com/in/yourprofile)

## Acknowledgments

- Built as a demonstration project for full-stack development skills
- Designed for use in technical interviews and portfolio presentations
- Implements industry-standard security and architectural patterns

## Support

For issues or questions, please open an issue on GitHub or contact the author.

---

**Note**: This project is intended for educational and portfolio purposes. For production use, additional security measures and optimizations are recommended.
