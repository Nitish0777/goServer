# Use a minimal golang image as a base
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download and cache Go dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin .

# Use a minimal alpine image as the final base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/bin .

# Expose port 8000 for the application
EXPOSE 8000

# Command to run the application
CMD ["./bin"]
