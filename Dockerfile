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

# Define an argument to accept the commit ID
ARG COMMIT_ID

# Set the commit ID as an environment variable
ENV COMMIT_ID=$COMMIT_ID

# Optionally, you can write the commit ID to a file for later use
RUN echo "Commit ID: $COMMIT_ID" > /etc/commit_id

# Copy the compiled Go binary from the builder
COPY --from=builder /app/nextcloud_uploader /app/nextcloud_uploader

# Make the binary executable
RUN chmod +x /app/nextcloud_uploader

# Verify that the binary exists
RUN ls -la /app/

# Define entrypoint to run the Go binary, and print the commit ID
ENTRYPOINT ["/bin/sh", "-c", "echo 'Commit ID: $COMMIT_ID'; cat /etc/commit_id; /app/nextcloud_uploader"]
