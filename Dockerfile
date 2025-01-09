# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the application source code
COPY . .

# Build the Go binary from cmd/main.go
RUN go build -o main ./cmd/main.go

# Stage 2: Create a minimal image with the built binary
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main ./

# Expose port 8080 for the application
EXPOSE 8080

# Run the binary
CMD ["./main"]