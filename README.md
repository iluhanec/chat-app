# Chat App

A simple Go-based chat application.

## Prerequisites

- Go 1.24+
- Docker (optional)

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

3. Install golangci-lint (for development):

   ```sh
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

4. Install goimports (for code formatting):

   ```sh
   go install golang.org/x/tools/cmd/goimports@latest
   ```

## Running & Building

This repo uses a **Makefile** for common tasks:

| Command           | Description                    |
| ----------------- | ------------------------------ |
| `make run`        | Run the application            |
| `make build`      | Compile the binary             |
| `make test`       | Execute all tests              |
| `make lint`       | Run code linter                |
| `make clean`      | Clean build artifacts          |
| `make docker-build` | Build Docker image            |
| `make docker-run` | Run Docker container          |

### Quick Start

```sh
make run
```

### Build Binary

```sh
make build
./chat-app
```

### Docker

```sh
make docker-build
make docker-run
```

## Development

- `main.go` - Application entry point
- `main_test.go` - Test suite
- `Makefile` - Build and development commands
- `.golangci.yml` - Linting configuration
- `Dockerfile` - Container configuration
- `.github/workflows/ci.yml` - CI/CD pipeline
