# Use the official Golang image for Go 1.22 on Alpine Linux 3.18
FROM golang:1.22.2-alpine3.18

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app for Linux (amd64 architecture)
RUN GOOS=linux GOARCH=amd64 go build -o main cmd/app/main.go

# Expose port 8080 to the outside world (if needed)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]