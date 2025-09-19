.PHONY: build test fmt run docker-up clean setup install-tools

BINARY_NAME=go-service
TOOLS=golang.org/x/tools/cmd/goimports golang.org/x/lint/golint github.com/evilmartians/lefthook

build:
	@echo "Building Go application..."
	go build -o $(BINARY_NAME) ./cmd/main.go
	@echo "Build successful."

test:
	@echo "Running tests..."
	go test ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

run: build
	@echo "Running application..."
	./$(BINARY_NAME)

docker-up:
	@echo "Starting Dockerized environment..."
	docker-compose up --build

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	docker-compose down

install-tools:
	@echo "Installing development tools..."
	go install $(TOOLS)

setup: install-tools
	@echo "Installing lefthook Git hooks..."
	@$(go env GOPATH)/bin/lefthook install
