package models

import "time"

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Chat represents a chat room
type Chat struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateChatRequest represents a request to create a new chat
type CreateChatRequest struct {
	Name string `json:"name"`
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
