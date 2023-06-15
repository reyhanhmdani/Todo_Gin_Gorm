FROM golang:1.20-alpine AS builder

LABEL author="Raihan hamdani" \
      title="todolist_api" \
      website="https://github.com/reyhanhmdani/Todo_Gin_Gorm"


RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o binary .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/binary .
COPY --from=builder /app/database/migrations/ ./database/migrations/

CMD ["./binary"]
