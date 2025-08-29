package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chat-app/internal/models"
)

func TestServer(t *testing.T) {
	server := NewServer()

	t.Run("ListChats_Empty", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/chats", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var chats []*models.Chat
		if err := json.NewDecoder(rr.Body).Decode(&chats); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if len(chats) != 0 {
			t.Errorf("Expected empty chat list, got %d chats", len(chats))
		}
	})

	t.Run("CreateChat", func(t *testing.T) {
		reqBody := models.CreateChatRequest{Name: "Test Chat"}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var chat models.Chat
		if err := json.NewDecoder(rr.Body).Decode(&chat); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if chat.Name != "Test Chat" {
			t.Errorf("Expected chat name 'Test Chat', got '%s'", chat.Name)
		}

		if chat.ID == "" {
			t.Error("Expected chat to have an ID")
		}
	})

	t.Run("CreateChat_EmptyName", func(t *testing.T) {
		reqBody := models.CreateChatRequest{Name: ""}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("SendAndGetMessages", func(t *testing.T) {
		// First create a chat
		createReq := models.CreateChatRequest{Name: "Message Test Chat"}
		createBody, _ := json.Marshal(createReq)

		req, _ := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		var chat models.Chat
		if err := json.NewDecoder(rr.Body).Decode(&chat); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		// Send a message
		msgReq := models.SendMessageRequest{
			Username: "testuser",
			Content:  "Hello, World!",
		}
		msgBody, _ := json.Marshal(msgReq)

		req, _ = http.NewRequest("POST", "/api/chats/"+chat.ID+"/messages", bytes.NewBuffer(msgBody))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var message models.Message
		if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
			t.Errorf("Failed to decode message response: %v", err)
		}

		if message.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", message.Username)
		}

		if message.Content != "Hello, World!" {
			t.Errorf("Expected content 'Hello, World!', got '%s'", message.Content)
		}

		// Get messages
		req, _ = http.NewRequest("GET", "/api/chats/"+chat.ID+"/messages", nil)
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var messages []*models.Message
		if err := json.NewDecoder(rr.Body).Decode(&messages); err != nil {
			t.Errorf("Failed to decode messages response: %v", err)
		}

		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}

		if messages[0].Content != "Hello, World!" {
			t.Errorf("Expected message content 'Hello, World!', got '%s'", messages[0].Content)
		}
	})

	t.Run("GetMessages_NonExistentChat", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/chats/nonexistent/messages", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("SendMessage_NonExistentChat", func(t *testing.T) {
		msgReq := models.SendMessageRequest{
			Username: "testuser",
			Content:  "Hello",
		}
		msgBody, _ := json.Marshal(msgReq)

		req, _ := http.NewRequest("POST", "/api/chats/nonexistent/messages", bytes.NewBuffer(msgBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("SendMessage_MissingFields", func(t *testing.T) {
		// First create a chat
		createReq := models.CreateChatRequest{Name: "Validation Test Chat"}
		createBody, _ := json.Marshal(createReq)

		req, _ := http.NewRequest("POST", "/api/chats", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		var chat models.Chat
		if err := json.NewDecoder(rr.Body).Decode(&chat); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		// Send message without username
		msgReq := models.SendMessageRequest{
			Username: "",
			Content:  "Hello",
		}
		msgBody, _ := json.Marshal(msgReq)

		req, _ = http.NewRequest("POST", "/api/chats/"+chat.ID+"/messages", bytes.NewBuffer(msgBody))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
