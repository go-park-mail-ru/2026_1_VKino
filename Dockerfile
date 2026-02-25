# Этап 1: Сборка (Builder)
FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /app

COPY go.mod  ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o myapp ./cmd/api

# Этап 2: Финальный образ (Runner)
FROM alpine:latest

# Устанавливаем сертификаты для работы HTTPS и таймзоны (опционально)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем скомпилированный бинарник из этапа builder
COPY --from=builder /app/myapp .

# Открываем порт 8080
EXPOSE 8080

CMD ["./myapp"]