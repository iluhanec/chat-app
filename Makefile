.PHONY: run-server run-client build test lint clean docker-build docker-run format format-check build-server build-client build-all

# Run the server
run-server:
	go run cmd/server/main.go

# Run the client
run-client:
	go run cmd/client/main.go --username=$(USER)

# Build server binary
build-server:
	go build -o bin/chat-server cmd/server/main.go

# Build client binary
build-client:
	go build -o bin/chat-client cmd/client/main.go

# Build client for all platforms
build-all: build-server
	@echo "Building client for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/chat-client-linux-amd64 cmd/client/main.go
	@echo "Building client for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/chat-client-darwin-amd64 cmd/client/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/chat-client-darwin-arm64 cmd/client/main.go
	@echo "Building client for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/chat-client-windows-amd64.exe cmd/client/main.go
	@echo "All builds complete!"

# Run tests
test:
	go test -v ./...

# Format code
format:
	goimports -w .

# Check formatting without changes
format-check:
	goimports -d .

# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image for server
docker-build:
	docker build -t chat-server .

# Run Docker container
docker-run:
	docker run -p 8080:8080 chat-server

# Run server and client for demo
demo: docker-build
	@echo "Starting server in Docker..."
	docker run -d --name chat-server-demo -p 8080:8080 chat-server
	@echo "Server started. You can now run 'make run-client' to connect."
	@echo "To stop the server, run: docker stop chat-server-demo && docker rm chat-server-demo"
