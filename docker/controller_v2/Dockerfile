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

# Install Air for live-reloading in development
RUN go install github.com/air-verse/air@latest
COPY air.toml .air.toml

# Build the application for production caching
RUN go build -o ./build/go-service ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Default command is to run Air for development.
# For production, this should be overridden to run the compiled binary.
CMD ["air", "-c", "air.toml"]