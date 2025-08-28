.PHONY: run build test lint

# Run the application
run:
	go run main.go

# Build the binary
build:
	go build -o chat-app main.go

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run