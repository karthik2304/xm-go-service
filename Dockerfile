# Dockerfile
FROM golang:1.23.4-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/run-service/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY .env .
# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Expose the application port
EXPOSE 9091

# Run the binary
CMD ["/app/main"]