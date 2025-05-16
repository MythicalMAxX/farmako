FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Create final image
FROM alpine:latest

WORKDIR /app

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy swagger documentation
COPY swagger/ /app/swagger/

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["./main"] 