FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

CMD ["./main"]
