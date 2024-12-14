# Use the official Golang image as the base image
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Set GOPATH and update PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Create a new stage for the final image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/main .

FROM gomicro/goose

# Copy the migrations folder into the container
ADD /sql/migrations/*.sql /migrations/
ADD entrypoint.sh /migrations/

# Set the entrypoint to the entrypoint script
ENTRYPOINT ["/migrations/entrypoint.sh"]

# Command to run the executable
CMD ["./main"]