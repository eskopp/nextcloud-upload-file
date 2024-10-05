# Use Go 1.23.2 to build the Go binary
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./go.mod
# COPY go.sum ./go.sum
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary (statically linked)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nextcloud_uploader .

# Use a lightweight image for running the binary
FROM alpine:latest

# Set working directory for runtime container
WORKDIR /app

# Install any necessary dependencies for alpine (if dynamic linking is required)
RUN apk add --no-cache libc6-compat

# Define build arguments and environment variables for the override flag
ARG OVERRIDE=false
ENV INPUT_OVERRIDE=$OVERRIDE

# Copy the compiled Go binary from the builder
COPY --from=builder /app/nextcloud_uploader /app/nextcloud_uploader

# Make the binary executable
RUN chmod +x /app/nextcloud_uploader

# Verify that the binary exists
RUN ls -la /app/

# Define entrypoint to run the Go binary and pass the override flag
ENTRYPOINT ["/app/nextcloud_uploader"]
