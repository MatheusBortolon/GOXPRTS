FROM golang:1.22-alpine AS base

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

FROM base AS test
RUN go test ./...

FROM base AS build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=build /bin/app /app

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app"]
