# Use the official Golang image to create a build artifact.
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app with CGO disabled to ensure static linking
RUN CGO_ENABLED=0 go build -o main cmd/main.go

# Use a smaller image for the final build
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the .env file
COPY .env .env

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]