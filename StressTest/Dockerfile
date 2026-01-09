FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go build -o stresstest .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/stresstest .

ENTRYPOINT ["./stresstest"]
