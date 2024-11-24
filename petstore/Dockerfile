# Use the official Golang image as the base
FROM golang:1.23-alpine AS builder

# Install make
RUN apk add --no-cache make

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN make build SERVICES=petstore

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/petstore/bin/main .

# Expose the port the app will run on
EXPOSE 8080

# Run the executable
CMD ["./main"]
