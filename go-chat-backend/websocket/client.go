package websocket

import (
	"context"
	"encoding/json"
	"log"

	"chat.app/backend/models"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		// Just forward raw bytes to Hub. Hub will unmarshal, timestamp, save, and broadcast.
		c.Hub.Broadcast <- messageBytes
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

// FetchHistory fetches chat history and sends it to a client
func (c *Client) FetchHistory() {
	if MessagesCollection == nil {
		return
	}

	// Get the RoomID from the client's hub
	roomID := c.Hub.RoomID

	// Set query options
	opts := options.Find()
	opts.SetLimit(50)
	opts.SetSort(bson.D{{Key: "timestamp", Value: 1}}) // Oldest first for chat history? No, usually newest last.
	// If we sort by timestamp 1, we get oldest first. That's what we want for appending to the list.

	// Filter by room_id (matching struct tag in models.Message)
	filter := bson.D{{Key: "room_id", Value: roomID}}

	cursor, err := MessagesCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Printf("error finding history: %v", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("error decoding history: %v", err)
			continue
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		c.Send <- msgBytes
	}
}
