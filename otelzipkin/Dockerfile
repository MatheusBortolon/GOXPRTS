FROM golang:1.21-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${SERVICE}

FROM alpine:3.19
WORKDIR /app
COPY --from=build /out/app /app/app
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/app/app"]
