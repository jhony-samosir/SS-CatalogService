# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/api/main.go

# Run stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Copy .env file if it exists
COPY .env .
# Copy migrations
COPY --from=builder /app/db ./db

EXPOSE 8081

# Command to run the executable
CMD ["./main"]
