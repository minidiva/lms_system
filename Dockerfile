# 1. Базовый образ
FROM golang:1.25-alpine AS builder

# 2. Установка зависимостей
RUN apk add git

# 3. Рабочая директория внутри контейнера
WORKDIR /app

# 4. Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# 5. Копируем остальной исходный код
COPY . .

# 6. Установка goose CLI (если надо использовать в миграциях)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# 7. Сборка приложения
RUN go build -o main ./cmd/lms

# Финальный контейнер
FROM alpine:latest

RUN apk add ca-certificates

WORKDIR /root/

# Копируем приложение и миграции
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# 🔧 Копируем goose CLI из builder
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Устанавливаем точку входа
ENTRYPOINT ["/entrypoint.sh"]