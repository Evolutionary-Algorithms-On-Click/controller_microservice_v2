# Use a recent Go version
FROM golang:1.25-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies. This layer is only re-built when go.mod/go.sum change.
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o ./build/go-service ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["go", "run", "-v", "./cmd/main.go"]