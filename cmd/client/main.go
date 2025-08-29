package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"chat-app/internal/models"
)

type Client struct {
	serverURL   string
	username    string
	currentChat string
	reader      *bufio.Reader
}

func NewClient(serverURL, username string) *Client {
	return &Client{
		serverURL: serverURL,
		username:  username,
		reader:    bufio.NewReader(os.Stdin),
	}
}

func (c *Client) Run() {
	fmt.Printf("Welcome to Chat App, %s!\n", c.username)
	fmt.Println("Commands:")
	fmt.Println("  /list        - List all chats")
	fmt.Println("  /create NAME - Create a new chat")
	fmt.Println("  /join ID     - Join a chat")
	fmt.Println("  /refresh     - Refresh messages in current chat")
	fmt.Println("  /quit        - Exit the application")
	fmt.Println()

	for {
		if c.currentChat != "" {
			fmt.Printf("[%s] > ", c.currentChat)
		} else {
			fmt.Print("> ")
		}

		input, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if strings.HasPrefix(input, "/") {
			c.handleCommand(input)
		} else if c.currentChat != "" {
			c.sendMessage(input)
		} else {
			fmt.Println("Please join a chat first using /join ID")
		}
	}
}

func (c *Client) handleCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/list":
		c.listChats()
	case "/create":
		if len(parts) < 2 {
			fmt.Println("Usage: /create NAME")
			return
		}
		name := strings.Join(parts[1:], " ")
		c.createChat(name)
	case "/join":
		if len(parts) != 2 {
			fmt.Println("Usage: /join ID")
			return
		}
		c.joinChat(parts[1])
	case "/refresh":
		if c.currentChat != "" {
			c.refreshMessages()
		} else {
			fmt.Println("Not in a chat")
		}
	case "/quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Unknown command:", parts[0])
	}
}

func (c *Client) listChats() {
	resp, err := http.Get(c.serverURL + "/api/chats")
	if err != nil {
		fmt.Println("Error fetching chats:", err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	var chats []*models.Chat
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if len(chats) == 0 {
		fmt.Println("No chats available. Create one with /create NAME")
		return
	}

	fmt.Println("\nAvailable chats:")
	for _, chat := range chats {
		fmt.Printf("  ID: %s | Name: %s | Created: %s\n",
			chat.ID[:8], chat.Name, chat.CreatedAt.Format("15:04:05"))
	}
	fmt.Println()
}

func (c *Client) createChat(name string) {
	reqBody, _ := json.Marshal(models.CreateChatRequest{Name: name})
	resp, err := http.Post(c.serverURL+"/api/chats", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating chat:", err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to create chat: %s\n", string(body))
		return
	}

	var chat models.Chat
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Printf("Created chat '%s' with ID: %s\n", chat.Name, chat.ID[:8])
	fmt.Println("Join it with: /join", chat.ID)
}

func (c *Client) joinChat(chatID string) {
	// First, check if chat exists by getting messages
	resp, err := http.Get(c.serverURL + "/api/chats/" + chatID + "/messages")
	if err != nil {
		fmt.Println("Error joining chat:", err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Chat not found")
		return
	}

	var messages []*models.Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		fmt.Println("Error decoding messages:", err)
		return
	}

	c.currentChat = chatID
	fmt.Printf("\nJoined chat %s\n", chatID[:8])
	fmt.Println("=== Chat History ===")

	if len(messages) == 0 {
		fmt.Println("(No messages yet)")
	} else {
		for _, msg := range messages {
			c.displayMessage(msg)
		}
	}
	fmt.Println("===================")
}

func (c *Client) refreshMessages() {
	resp, err := http.Get(c.serverURL + "/api/chats/" + c.currentChat + "/messages")
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	var messages []*models.Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		fmt.Println("Error decoding messages:", err)
		return
	}

	fmt.Println("\n=== Refreshed Messages ===")
	if len(messages) == 0 {
		fmt.Println("(No messages yet)")
	} else {
		for _, msg := range messages {
			c.displayMessage(msg)
		}
	}
	fmt.Println("========================")
}

func (c *Client) sendMessage(content string) {
	reqBody, _ := json.Marshal(models.SendMessageRequest{
		Username: c.username,
		Content:  content,
	})

	resp, err := http.Post(
		c.serverURL+"/api/chats/"+c.currentChat+"/messages",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to send message: %s\n", string(body))
	}
}

func (c *Client) displayMessage(msg *models.Message) {
	timestamp := msg.Timestamp.Format("15:04:05")
	if msg.Username == c.username {
		fmt.Printf("[%s] You: %s\n", timestamp, msg.Content)
	} else {
		fmt.Printf("[%s] %s: %s\n", timestamp, msg.Username, msg.Content)
	}
}

func main() {
	username := flag.String("username", "", "Your username")
	server := flag.String("server", "http://localhost:8080", "Server URL")
	flag.Parse()

	if *username == "" {
		fmt.Println("Error: --username flag is required")
		flag.Usage()
		os.Exit(1)
	}

	client := NewClient(*server, *username)
	client.Run()
}
