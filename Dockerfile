# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git build-base alsa-lib-dev

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

# Install runtime dependencies
RUN apk add --no-cache ca-certificates alsa-lib

# Copy binary from builder
COPY --from=builder /app/simple-google-speaker .

# Create audio directory
RUN mkdir -p audio

# Expose port (default 8080)
EXPOSE 8080

# Environment variables with defaults
ENV PORT=:8080
ENV AUDIO_FOLDER=audio
ENV VOLUME=100

# Run the application
CMD ["./simple-google-speaker"]
