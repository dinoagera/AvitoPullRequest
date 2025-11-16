FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o /app/main ./cmd/main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.18

COPY --from=builder /app/main /main
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations

COPY .env .env

EXPOSE 8080

CMD ["sh", "-c", "goose -dir ./migrations postgres \"$DB_URL\" up && ./main"]