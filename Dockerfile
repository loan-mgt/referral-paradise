# Multi-stage build

# Stage 1: Build the Go binary
FROM golang:1.23.0-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create a minimal image to run the Go application
FROM alpine:3.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built Go binary from the build stage
COPY --from=build /app/main /app/main

# Expose the application on port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
