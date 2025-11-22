package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"chat.app/backend/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var Mongo *mongo.Client
var MessagesCollection *mongo.Collection

type Hub struct {
	RoomID     string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub(roomID string) *Hub {
	return &Hub{
		RoomID:     roomID,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case messageBytes := <-h.Broadcast:
			// 1. Unmarshal
			var msg models.Message
			if err := json.Unmarshal(messageBytes, &msg); err != nil {
				log.Println("Error unmarshaling message:", err)
				continue
			}

			// 2. Set server-side timestamp and roomID
			if msg.Timestamp.IsZero() {
				msg.Timestamp = time.Now()
			}
			if msg.RoomID == "" {
				msg.RoomID = h.RoomID
			}

			// 3. Save to MongoDB
			h.saveMessage(msg)

			// 4. Marshal back to JSON
			finalBytes, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshaling message:", err)
				continue
			}

			// 5. Broadcast to clients
			for client := range h.Clients {
				select {
				case client.Send <- finalBytes:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) saveMessage(msg models.Message) {
	if MessagesCollection == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := MessagesCollection.InsertOne(ctx, msg)
	if err != nil {
		log.Println("Error saving message to DB:", err)
	}
}

type HubManager struct {
	Hubs map[string]*Hub
	mu   sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		Hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) GetHub(roomID string) *Hub {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	if hub, ok := hm.Hubs[roomID]; ok {
		return hub
	}

	hub := NewHub(roomID)
	hm.Hubs[roomID] = hub
	go hub.Run()
	return hub
}

func (hm *HubManager) ListRooms() []string {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	rooms := make([]string, 0, len(hm.Hubs))
	for roomID := range hm.Hubs {
		rooms = append(rooms, roomID)
	}
	return rooms
}

// CloseHub closes a hub and disconnects all its clients
func (hm *HubManager) CloseHub(roomID string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hub, exists := hm.Hubs[roomID]
	if !exists {
		return
	}

	// Close all client connections in this hub
	for client := range hub.Clients {
		close(client.Send)
		client.Conn.Close()
		delete(hub.Clients, client)
	}

	// Remove hub from manager
	delete(hm.Hubs, roomID)
	log.Printf("Hub for room '%s' has been closed", roomID)
}
