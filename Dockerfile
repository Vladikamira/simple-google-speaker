# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o simple-google-speaker main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (like ca-certificates for HTTPS requests to TTS)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/simple-google-speaker .

# Create audio directory
RUN mkdir -p audio

# Expose port (default 8080)
EXPOSE 8080

# Environment variables with defaults
ENV PORT=:8080
ENV AUDIO_FOLDER=audio
ENV LANGUAGE=en
ENV MESSAGE_TEXT="Time to sleep"
ENV VOLUME=100

# Run the application
CMD ["./simple-google-speaker"]
