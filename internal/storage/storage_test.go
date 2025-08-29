package storage

import (
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	t.Run("CreateChat", func(t *testing.T) {
		s := NewStorage()
		chat, err := s.CreateChat("Test Chat")

		if err != nil {
			t.Errorf("Failed to create chat: %v", err)
		}

		if chat.Name != "Test Chat" {
			t.Errorf("Expected chat name 'Test Chat', got '%s'", chat.Name)
		}

		if chat.ID == "" {
			t.Error("Expected chat to have an ID")
		}

		if chat.CreatedAt.IsZero() {
			t.Error("Expected chat to have a creation time")
		}
	})

	t.Run("GetChat", func(t *testing.T) {
		s := NewStorage()
		created, _ := s.CreateChat("Test Chat")

		retrieved, exists := s.GetChat(created.ID)
		if !exists {
			t.Error("Expected chat to exist")
		}

		if retrieved.ID != created.ID {
			t.Errorf("Expected chat ID %s, got %s", created.ID, retrieved.ID)
		}

		if retrieved.Name != created.Name {
			t.Errorf("Expected chat name %s, got %s", created.Name, retrieved.Name)
		}
	})

	t.Run("GetChat_NonExistent", func(t *testing.T) {
		s := NewStorage()
		_, exists := s.GetChat("nonexistent")
		if exists {
			t.Error("Expected chat to not exist")
		}
	})

	t.Run("ListChats", func(t *testing.T) {
		s := NewStorage()

		// Initially empty
		chats := s.ListChats()
		if len(chats) != 0 {
			t.Errorf("Expected 0 chats, got %d", len(chats))
		}

		// Create some chats
		if _, err := s.CreateChat("Chat 1"); err != nil {
			t.Errorf("Failed to create chat 1: %v", err)
		}
		if _, err := s.CreateChat("Chat 2"); err != nil {
			t.Errorf("Failed to create chat 2: %v", err)
		}
		if _, err := s.CreateChat("Chat 3"); err != nil {
			t.Errorf("Failed to create chat 3: %v", err)
		}

		chats = s.ListChats()
		if len(chats) != 3 {
			t.Errorf("Expected 3 chats, got %d", len(chats))
		}
	})

	t.Run("AddMessage", func(t *testing.T) {
		s := NewStorage()
		chat, _ := s.CreateChat("Test Chat")

		message, err := s.AddMessage(chat.ID, "testuser", "Hello, World!")
		if err != nil {
			t.Errorf("Failed to add message: %v", err)
		}

		if message == nil {
			t.Fatal("Expected message to be returned")
		}

		if message.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", message.Username)
		}

		if message.Content != "Hello, World!" {
			t.Errorf("Expected content 'Hello, World!', got '%s'", message.Content)
		}

		if message.ChatID != chat.ID {
			t.Errorf("Expected chat ID %s, got %s", chat.ID, message.ChatID)
		}

		if message.Timestamp.IsZero() {
			t.Error("Expected message to have a timestamp")
		}
	})

	t.Run("AddMessage_NonExistentChat", func(t *testing.T) {
		s := NewStorage()
		message, err := s.AddMessage("nonexistent", "testuser", "Hello")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if message != nil {
			t.Error("Expected nil message for non-existent chat")
		}
	})

	t.Run("GetMessages", func(t *testing.T) {
		s := NewStorage()
		chat, _ := s.CreateChat("Test Chat")

		// Initially empty
		messages, exists := s.GetMessages(chat.ID)
		if !exists {
			t.Error("Expected messages array to exist for created chat")
		}

		if len(messages) != 0 {
			t.Errorf("Expected 0 messages, got %d", len(messages))
		}

		// Add some messages
		if _, err := s.AddMessage(chat.ID, "user1", "Message 1"); err != nil {
			t.Errorf("Failed to add message 1: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
		if _, err := s.AddMessage(chat.ID, "user2", "Message 2"); err != nil {
			t.Errorf("Failed to add message 2: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		if _, err := s.AddMessage(chat.ID, "user1", "Message 3"); err != nil {
			t.Errorf("Failed to add message 3: %v", err)
		}

		messages, exists = s.GetMessages(chat.ID)
		if !exists {
			t.Error("Expected messages to exist")
		}

		if len(messages) != 3 {
			t.Errorf("Expected 3 messages, got %d", len(messages))
		}

		// Verify message order (should be in chronological order)
		if messages[0].Content != "Message 1" {
			t.Errorf("Expected first message to be 'Message 1', got '%s'", messages[0].Content)
		}

		if messages[2].Content != "Message 3" {
			t.Errorf("Expected last message to be 'Message 3', got '%s'", messages[2].Content)
		}
	})

	t.Run("GetMessages_NonExistentChat", func(t *testing.T) {
		s := NewStorage()
		_, exists := s.GetMessages("nonexistent")
		if exists {
			t.Error("Expected messages to not exist for non-existent chat")
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		s := NewStorage()
		chat, _ := s.CreateChat("Concurrent Chat")

		// Test concurrent writes
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(n int) {
				if _, err := s.AddMessage(chat.ID, "user", "Message"); err != nil {
					t.Errorf("Failed to add message in goroutine %d: %v", n, err)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		messages, _ := s.GetMessages(chat.ID)
		if len(messages) != 10 {
			t.Errorf("Expected 10 messages after concurrent writes, got %d", len(messages))
		}
	})
}
