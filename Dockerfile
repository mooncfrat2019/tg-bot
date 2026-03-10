# Dockerfile
# Этап 1: Сборка Go приложений
FROM golang:1.22-alpine AS builder

# Устанавливаем необходимые системные зависимости
RUN apk add --no-cache git gcc musl-dev make

# Устанавливаем рабочую директорию
WORKDIR /build

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем оба приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./ && \
    CGO_ENABLED=0 GOOS=linux go build -o cron ./cmd/cron

# Этап 2: Финальный образ
FROM node:22-alpine

# Устанавливаем PM2 глобально
RUN npm install -g pm2 --silent

# Устанавливаем CA сертификаты и tzdata
RUN apk --no-cache add ca-certificates tzdata

# Создаем непривилегированного пользователя
#RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Создаем рабочую директорию
WORKDIR /app
RUN chown -R node:node /app

# Копируем бинарные файлы из builder
COPY --from=builder /build/app /app/app
COPY --from=builder /build/cron /app/cron

# Копируем конфигурацию PM2
COPY ecosystem.config.cjs /app/

# Копируем скрипт запуска
COPY start.sh /app/

# Делаем бинарники исполняемыми и устанавливаем права
RUN chmod +x /app/app && \
    chmod +x /app/cron && \
    chmod +x /app/start.sh && \
    chown -R node:node /app
# Проверяем, что файлы скопировались
RUN ls -la /app/

USER node
# Объявляем порт для веб приложения
EXPOSE 8765
# Запускаем через start.sh
CMD ["/app/start.sh"]