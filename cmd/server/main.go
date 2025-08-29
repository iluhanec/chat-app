package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"chat-app/internal/models"
	"chat-app/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	storage *storage.Storage
	router  *mux.Router
}

func NewServer() *Server {
	s := &Server{
		storage: storage.NewStorage(),
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/api/chats", s.handleListChats).Methods("GET")
	s.router.HandleFunc("/api/chats", s.handleCreateChat).Methods("POST")
	s.router.HandleFunc("/api/chats/{chatID}/messages", s.handleGetMessages).Methods("GET")
	s.router.HandleFunc("/api/chats/{chatID}/messages", s.handleSendMessage).Methods("POST")
}

func (s *Server) handleListChats(w http.ResponseWriter, r *http.Request) {
	chats := s.storage.ListChats()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCreateChat(w http.ResponseWriter, r *http.Request) {
	var req models.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Chat name is required", http.StatusBadRequest)
		return
	}

	chat, err := s.storage.CreateChat(req.Name)
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatID"]

	messages, exists := s.storage.GetMessages(chatID)
	if !exists {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatID"]

	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Content == "" {
		http.Error(w, "Username and content are required", http.StatusBadRequest)
		return
	}

	message, err := s.storage.AddMessage(chatID, req.Username, req.Content)
	if err != nil || message == nil {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	server := NewServer()

	log.Printf("Starting server on port %s", *port)

	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      server.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
