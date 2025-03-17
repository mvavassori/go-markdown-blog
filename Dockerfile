# Start from the official Golang base image
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o go-markdown-blog .

# Start a new stage from a smaller base image
FROM debian:stable-slim

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/go-markdown-blog .

# Copy the templates, static, and posts directories
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/posts ./posts

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./go-markdown-blog"]