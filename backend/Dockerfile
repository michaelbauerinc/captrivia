# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

RUN go test -v .

# Build the Go app - specify the new location of the main.go file
RUN go build -o main ./cmd


# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
