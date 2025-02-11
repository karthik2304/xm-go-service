FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/run-service/main.go

FROM alpine:latest

WORKDIR /app

COPY .env .
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 9091

CMD ["/app/main"]