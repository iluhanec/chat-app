.PHONY: run build test lint clean docker-build docker-run format format-check

# Run the application
run:
	go run main.go

# Build the binary
build:
	go build -o chat-app main.go

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
	rm -f chat-app
	go clean

# Build Docker image
docker-build:
	docker build -t chat-app .

# Run Docker container
docker-run:
	docker run -p 8080:8080 chat-app