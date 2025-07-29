# Use official Golang image that matches go.mod version
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Copy rest of the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port
EXPOSE 8080

# Start the app
CMD ["./main"]
