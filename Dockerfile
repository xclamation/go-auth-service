FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM gomicro/goose
WORKDIR /app
COPY /sql/migrations/*.sql /app/migrations/
COPY entrypoint.sh /app/migrations/
RUN chmod +x /app/migrations/entrypoint.sh
COPY --from=builder /app/main .
RUN chmod +x ./main
COPY .env .env

ENTRYPOINT ["/app/migrations/entrypoint.sh"]