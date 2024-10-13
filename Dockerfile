# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-cli ./cmd/main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Install MySQL client for database connections if needed
RUN apk --no-cache add mysql-client

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/todo-cli .

# Expose the gRPC port
EXPOSE 50054

# Start the gRPC server
CMD ["./todo-cli", "--db-user=root", "--db-password=my_secure_password", "--db-host=db", "--db-port=3306", "--db-name=todo_app", "--grpc-port=50054"]
