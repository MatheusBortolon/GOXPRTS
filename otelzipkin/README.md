# Otel + Zipkin - Temperature by CEP

This project implements two services:
- Service A: receives CEP input and forwards to Service B
- Service B: resolves CEP to city and returns temperature in Celsius, Fahrenheit, and Kelvin

It includes distributed tracing with OpenTelemetry and Zipkin.

## Requirements

- Docker and Docker Compose
- WeatherAPI key

## Setup

1. Copy the environment file:

```bash
cp .env.example .env
```

2. Add your WeatherAPI key in `.env`.

## Run with Docker Compose

```bash
docker compose up --build
```

Service A runs on http://localhost:8080
Zipkin runs on http://localhost:9411

## Usage

### Valid request

```bash
curl -X POST http://localhost:8080/zipcode \
  -H "Content-Type: application/json" \
  -d '{"cep":"29902555"}'
```

### Invalid CEP

```bash
curl -X POST http://localhost:8080/zipcode \
  -H "Content-Type: application/json" \
  -d '{"cep":"123"}'
```

### Expected responses

- 200: `{ "city": "Sao Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }`
- 422: `invalid zipcode`
- 404: `can not find zipcode`

## Tracing

- Service A and Service B export traces to the OpenTelemetry Collector.
- The collector exports to Zipkin.
- Open Zipkin UI at http://localhost:9411

## Local run (optional)

```bash
go mod download
go run ./cmd/service-b
```

```bash
SERVICE_B_URL=http://localhost:8081 go run ./cmd/service-a
```
