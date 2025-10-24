# Multi-stage Dockerfile for building and running the Go application
# This demonstrates best practices for containerizing Go applications

# Stage 1: Build the application
# We use the official Go image as our build environment
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates which are needed for go mod download
# and for HTTPS connections
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the go module files first
# This is done separately from copying the source code because Docker caches
# layers. If your go.mod and go.sum don't change, Docker can reuse this layer
# even if your source code changes, which speeds up builds.
COPY go.mod go.sum* ./

# Download dependencies
# go.sum might not exist yet if you haven't run go mod tidy, so we made it optional above
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary that doesn't depend on libc
# This allows us to use a minimal base image in the next stage
# The -ldflags="-w -s" flags strip debug information to make the binary smaller
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .

# Stage 2: Create the runtime image
# We use alpine as our base image because it's tiny (about 5MB)
# This keeps our final image small, which is faster to deploy and more secure
FROM alpine:latest

# Install ca-certificates so our app can make HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Create a non-root user to run the application
# Running as root in containers is a security risk
RUN adduser -D -u 1000 appuser

# Set the working directory
WORKDIR /home/appuser

# Copy the binary from the builder stage
# We're copying from the builder stage (named above) to this stage
COPY --from=builder /app/server .

# Change ownership of the binary to our non-root user
RUN chown appuser:appuser server

# Switch to the non-root user
USER appuser

# Document which port the application listens on
# This doesn't actually publish the port, it's documentation for developers
# and tools that inspect the image
EXPOSE 8000

# Set the entrypoint to run our binary
# Using ENTRYPOINT instead of CMD means this command always runs,
# and any arguments passed to 'docker run' will be passed to our binary
ENTRYPOINT ["./server"]
