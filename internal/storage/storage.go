package storage

import (
	"sync"
	"time"

	"chat-app/internal/models"

	"github.com/google/uuid"
)

// Storage is an in-memory storage for chats and messages
type Storage struct {
	mu       sync.RWMutex
	chats    map[string]*models.Chat
	messages map[string][]*models.Message // chatID -> messages
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	return &Storage{
		chats:    make(map[string]*models.Chat),
		messages: make(map[string][]*models.Message),
	}
}

// CreateChat creates a new chat
func (s *Storage) CreateChat(name string) (*models.Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat := &models.Chat{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	s.chats[chat.ID] = chat
	s.messages[chat.ID] = []*models.Message{}

	return chat, nil
}

// GetChat retrieves a chat by ID
func (s *Storage) GetChat(chatID string) (*models.Chat, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, exists := s.chats[chatID]
	return chat, exists
}

// ListChats returns all chats
func (s *Storage) ListChats() []*models.Chat {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chats := make([]*models.Chat, 0, len(s.chats))
	for _, chat := range s.chats {
		chats = append(chats, chat)
	}
	return chats
}

// AddMessage adds a message to a chat
func (s *Storage) AddMessage(chatID, username, content string) (*models.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.chats[chatID]; !exists {
		return nil, nil
	}

	message := &models.Message{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}

	s.messages[chatID] = append(s.messages[chatID], message)
	return message, nil
}

// GetMessages retrieves all messages for a chat
func (s *Storage) GetMessages(chatID string) ([]*models.Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages, exists := s.messages[chatID]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	result := make([]*models.Message, len(messages))
	copy(result, messages)
	return result, true
}
