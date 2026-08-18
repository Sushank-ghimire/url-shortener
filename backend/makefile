APP_NAME=url-shortener
MAIN_PATH=./cmd/api

.PHONY: run build test clean tidy migrate-up migrate-down

# Run the application
run:
	go run $(MAIN_PATH)

# Run the build server
serve: build
	./bin/$(APP_NAME)

# Build binary
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

# Run tests
test:
	go test ./...

# Run tests with coverage
coverage:
	go test -cover ./...

# Clean build files
clean:
	rm -rf bin/

# Update dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Check code issues
vet:
	go vet ./...

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down
