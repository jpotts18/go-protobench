package model

import "time"

// Message represents the common message structure used across all protocols
type Message struct {
	ID        string    `json:"id" bson:"id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Content   string    `json:"content" bson:"content"`
	Number    int64     `json:"number" bson:"number"`
	IsValid   bool      `json:"is_valid" bson:"is_valid"`
}

// Protocol defines the interface that all protocol implementations must satisfy
type Protocol interface {
	Name() string
	StartServer() error
	StopServer() error
	SendMessage(msg *Message) error
}
