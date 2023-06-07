FROM golang:1.20-alpine AS builder


LABEL author="Raihan hamdani"
LABEL Title="todolist_api" website="https://github.com/reyhanhmdani/Todo_Gin_Gorm"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o binary ./

FROM scratch

COPY --from=builder /app/binary .
COPY --from=builder /app/database/migrations/ ./app/database/migrations/

WORKDIR /app

CMD ["/binary"]

# sudo docker run --network=host todolist
# docker rm -f $(docker ps -aq)
# docker network inspect my-network
