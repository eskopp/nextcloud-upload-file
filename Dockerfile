# Use Go 1.23.2 to build the Go binary
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o nextcloud_uploader .

# Use a lightweight image for running the binary
FROM alpine:latest

# Set working directory for runtime container
WORKDIR /app

# Copy the compiled Go binary from the builder
COPY --from=builder /app/nextcloud_uploader /app/nextcloud_uploader

# Make the binary executable
RUN chmod +x /app/nextcloud_uploader

# Define entrypoint to run the Go binary
ENTRYPOINT ["/app/nextcloud_uploader"]
