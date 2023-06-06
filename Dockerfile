FROM golang:1.20-alpine AS builder

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

