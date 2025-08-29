# Chat App

A minimalistic chat application with HTTP REST API server and console client.

## Features

- **Chat Server**: HTTP REST API for chat operations
- **Console Client**: Interactive terminal-based chat client
- **Multiple Chat Rooms**: Create and join different chat rooms
- **Chat History**: View all messages when joining a chat
- **Cross-Platform**: Client binaries for Linux, macOS, and Windows
- **Dockerized Server**: Easy deployment with Docker

## Architecture

```
chat-app/
├── cmd/
│   ├── server/        # HTTP server application
│   └── client/        # Console client application
├── internal/
│   ├── models/        # Data models
│   └── storage/       # In-memory storage
└── bin/               # Compiled binaries
```

## Prerequisites

- Go 1.21+
- Docker (optional, for containerized server)

## Quick Start

### Using Docker (Recommended)

1. Start the server:
   ```sh
   make demo
   ```

2. In another terminal, connect with the client:
   ```sh
   make run-client
   ```

### Manual Setup

1. Start the server:
   ```sh
   make run-server
   ```

2. In another terminal, connect with the client:
   ```sh
   go run cmd/client/main.go --username="YourName"
   ```

## Client Commands

Once connected, use these commands:

- `/list` - List all available chats
- `/create NAME` - Create a new chat room
- `/join ID` - Join an existing chat by ID
- `/refresh` - Refresh messages in current chat
- `/quit` - Exit the application

Any text without a `/` prefix will be sent as a message to the current chat.

## Building

### Build Everything

```sh
make build-all
```

This creates binaries for:
- Server: `bin/chat-server`
- Client:
  - Linux: `bin/chat-client-linux-amd64`
  - macOS Intel: `bin/chat-client-darwin-amd64`
  - macOS Apple Silicon: `bin/chat-client-darwin-arm64`
  - Windows: `bin/chat-client-windows-amd64.exe`

### Build Specific Components

```sh
make build-server     # Build server only
make build-client     # Build client for current platform
```

## API Endpoints

The server exposes the following REST API:

- `GET /api/chats` - List all chats
- `POST /api/chats` - Create a new chat
- `GET /api/chats/{chatID}/messages` - Get messages for a chat
- `POST /api/chats/{chatID}/messages` - Send a message to a chat

## Testing

Run comprehensive tests:

```sh
make test
```

## Development

### Available Make Commands

| Command             | Description                     |
| ------------------- | ------------------------------- |
| `make run-server`   | Run the server locally          |
| `make run-client`   | Run the client locally          |
| `make build-all`    | Build all platform binaries     |
| `make test`         | Run all tests                   |
| `make format`       | Format code with goimports      |
| `make lint`         | Run golangci-lint               |
| `make docker-build` | Build Docker image              |
| `make docker-run`   | Run server in Docker            |
| `make demo`         | Start server in Docker for demo |
| `make clean`        | Clean build artifacts           |

### Project Structure

- `cmd/server/` - HTTP server implementation
- `cmd/client/` - Console client implementation
- `internal/models/` - Shared data structures
- `internal/storage/` - In-memory storage layer

## Example Usage

1. Start the server:
   ```sh
   make docker-run
   ```

2. Connect first user:
   ```sh
   ./bin/chat-client --username="Alice"
   > /create General
   > /join <chat-id>
   > Hello everyone!
   ```

3. Connect second user:
   ```sh
   ./bin/chat-client --username="Bob"
   > /list
   > /join <chat-id>
   > Hi Alice!
   ```

## Notes

- No authentication/authorization (as per requirements)
- Messages are stored in-memory (lost on server restart)
- Server runs on port 8080 by default
- Client connects to http://localhost:8080 by default
