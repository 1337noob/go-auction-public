FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /auction ./cmd/main.go

FROM alpine:latest

# Устанавливаем tzdata для поддержки часовых поясов
RUN apk add --no-cache tzdata

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /auction /auction

EXPOSE 8081

CMD ["/auction"]
