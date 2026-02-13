# labGCR

Serviço em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## Requisitos

- Go 1.22+
- Chave da WeatherAPI (env `WEATHER_API_KEY`)

## Executar localmente

```bash
go test ./...
```

```bash
WEATHER_API_KEY=SUACHAVE go run ./cmd/server
```

Endpoint:

- `GET /weather?cep=01001000`

Resposta de sucesso (200):

```json
{ "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }
```

Erros:

- 422 `invalid zipcode`
- 404 `can not find zipcode`

## Docker

```bash
docker compose up --build
```

```bash
docker compose run --rm test
```

## Variaveis de ambiente

- `WEATHER_API_KEY` (obrigatoria)
- `PORT` (default: 8080)
- `CEP_API_BASE` (default: https://viacep.com.br)
- `WEATHER_API_BASE` (default: https://api.weatherapi.com)
