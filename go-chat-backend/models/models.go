package models

import "time"

type User struct {
	Username     string `json:"username" bson:"username"`
	Password     string `json:"password_hash" bson:"-"` // For receiving password from frontend
	PasswordHash string `json:"-" bson:"password_hash"` // For storing hash in DB
}

type Room struct {
	Name      string    `json:"name" bson:"name"`
	Creator   string    `json:"creator" bson:"creator"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Message struct {
	RoomID    string    `json:"room_id" bson:"room_id"`
	Username  string    `json:"username" bson:"username"`
	Content   string    `json:"content" bson:"content"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
