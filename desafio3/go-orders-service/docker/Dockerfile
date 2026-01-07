FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rest ./cmd/rest
RUN go build -o grpc ./cmd/grpc
RUN go build -o graphql ./cmd/graphql

FROM gcr.io/distroless/base

COPY --from=builder /app/rest /usr/local/bin/rest
COPY --from=builder /app/grpc /usr/local/bin/grpc
COPY --from=builder /app/graphql /usr/local/bin/graphql

CMD ["/usr/local/bin/rest"]