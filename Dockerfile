FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM gomicro/goose AS migrator
WORKDIR /app
#COPY --from=builder /app/main /app/main
RUN mkdir -p /app/migrations
COPY /sql/migrations/*.sql /app/migrations/
COPY entrypoint.sh /app/migrations/
RUN chmod +x /app/migrations/entrypoint.sh

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=migrator /app/migrations /app/migrations

ENTRYPOINT ["/bin/sh", "/app/migrations/entrypoint.sh"]

#CMD ["exec","./main"]