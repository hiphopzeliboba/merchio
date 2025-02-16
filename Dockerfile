FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Финальный этап
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]