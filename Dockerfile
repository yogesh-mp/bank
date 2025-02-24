# Use official Golang image as a build stage
FROM golang:1.23 as builder

# Set working directory inside container
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api/main.go

# Use a minimal image for running the binary
FROM alpine:latest

# Set working directory inside container
WORKDIR /app

# Install dependencies
RUN apk add --no-cache ca-certificates

# Copy built binary from the builder stage
COPY --from=builder /app/api ./

# Copy .env file
COPY .env .env

# Expose API server port
EXPOSE 8080

# Start API
CMD ["./api"]
