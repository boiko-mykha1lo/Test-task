FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notifications ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/notifications .

CMD ["./notifications"]