package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"chat.app/backend/models"
	"chat.app/backend/websocket"

	ws "github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// Global DB collections
var usersCollection *mongo.Collection
var roomsCollection *mongo.Collection

// serveWs handles websocket requests from clients
func serveWs(w http.ResponseWriter, r *http.Request, hubManager *websocket.HubManager) {
	// Extract room ID from URL
	roomID := r.URL.Path[len("/ws/"):]
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	// Extract username from query params (sent from frontend)
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous" // Fallback if not provided
	}

	// Ensure room exists in DB
	ensureRoomExists(roomID, username)

	// Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &websocket.Client{
		Hub:  hubManager.GetHub(roomID),
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
	go client.FetchHistory()
}

func ensureRoomExists(name, creator string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	var room models.Room
	err := roomsCollection.FindOne(ctx, filter).Decode(&room)
	if err == mongo.ErrNoDocuments {
		newRoom := models.Room{
			Name:      name,
			Creator:   creator,
			CreatedAt: time.Now(),
		}
		roomsCollection.InsertOne(ctx, newRoom)
		fmt.Printf("Created new room in DB: %s (by %s)\n", name, creator)
	}
}

func initDB() {
	connectionString := os.Getenv("MONGO_URI")
	if connectionString == "" {
		// Fallback for local development (replace with your local MongoDB)
		connectionString = "mongodb://localhost:27017"
		log.Println("Warning: MONGO_URI not set, using local MongoDB")
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	db := client.Database("chatapp")
	usersCollection = db.Collection("users")
	roomsCollection = db.Collection("rooms")

	// Pass collections to websocket package
	websocket.Mongo = client
	websocket.MessagesCollection = db.Collection("messages")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Register decode error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Registering user: %s", user.Username)

	// Hash the plain password from user.Password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Hash error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(hash)
	user.Password = "" // Clear the plaintext password

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Insert error (duplicate?):", err)
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	log.Println("User registered successfully")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password_hash"` // Frontend sends password_hash
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Println("Login decode error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Attempting login for: %s", creds.Username)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		log.Println("User not found:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
	if err != nil {
		log.Println("Password mismatch:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Println("Login successful")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "username": user.Username})
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	initDB()

	hubManager := websocket.NewHubManager()

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, hubManager)
	})

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			return
		}

		if r.Method == http.MethodDelete {
			// Delete room - extract room name from query params
			roomName := r.URL.Query().Get("name")
			username := r.URL.Query().Get("username")

			if roomName == "" || username == "" {
				http.Error(w, "Room name and username required", http.StatusBadRequest)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Check if user is the creator
			var room models.Room
			err := roomsCollection.FindOne(ctx, bson.M{"name": roomName}).Decode(&room)
			if err != nil {
				http.Error(w, "Room not found", http.StatusNotFound)
				return
			}

			if room.Creator != username {
				http.Error(w, "Only the room creator can delete this room", http.StatusForbidden)
				return
			}

			// Delete room from database
			_, err = roomsCollection.DeleteOne(ctx, bson.M{"name": roomName})
			if err != nil {
				http.Error(w, "Failed to delete room", http.StatusInternalServerError)
				return
			}

			// Close all WebSocket connections in that room
			hubManager.CloseHub(roomName)

			log.Printf("Room '%s' deleted by %s", roomName, username)
			json.NewEncoder(w).Encode(map[string]string{"message": "Room deleted successfully"})
			return
		}

		// GET - List all rooms
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := roomsCollection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var rooms []models.Room
		if err = cursor.All(ctx, &rooms); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return full room objects (including creator)
		json.NewEncoder(w).Encode(rooms)
	})

	fmt.Println("Server starting on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
