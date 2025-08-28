# Chat App

A simple Go-based chat application.

## Prerequisites

- Go 1.20+

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/chat-app.git
   cd chat-app
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

## Running & Building

This repo uses a **Makefile** for common tasks:

| Command      | Description         |
| ------------ | ------------------- |
| `make run`   | Run the application |
| `make build` | Compile the binary  |
| `make test`  | Execute all tests   |
| `make lint`  | Run code linter     |

### Quick Start

```sh
make run
```

### Build Binary

```sh
make build
./chat-app
```

Or directly with Go:

```sh
go build -o chat-app main.go
./chat-app
```
