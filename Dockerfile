# Используем официальный образ Go для сборки
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы модулей Go
COPY go.mod ./
#go.sum ./
# Скачиваем все зависимости
RUN go mod download

# Копируем весь остальной код в контейнер
COPY . .

# Собираем наше приложение. Бинарный файл будет называться 'gonoty'
RUN go build -o gonoty ./cmd/main.go

# Финальный этап: создаем легкий образ для запуска
FROM alpine:latest

# Устанавливаем необходимые библиотеки (для работы с SSL и т.д.)
RUN apk --no-cache add ca-certificates

# Копируем скомпилированный бинарник из этапа 'builder'
COPY --from=builder /app/gonoty .

# Указываем команду, которая выполнится при запуске контейнера
CMD ["./gonoty"]