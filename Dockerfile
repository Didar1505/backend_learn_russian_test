# ===== 1 stage: build =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum отдельно (для кэша)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api

# ===== 2 stage: runtime =====
FROM alpine:3.20

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]