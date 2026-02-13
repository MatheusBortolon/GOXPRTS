# Rate Limiter em Go

Sistema de rate limiting em Go que controla o tráfego de requisições HTTP com base em endereço IP ou token de acesso. A solução utiliza Redis para armazenamento e oferece configuração flexível através de variáveis de ambiente.

## 🎯 Características

- ✅ Limitação por endereço IP
- ✅ Limitação por token de acesso (API_KEY)
- ✅ Priorização de limites por token sobre IP
- ✅ Middleware HTTP reutilizável
- ✅ Armazenamento em Redis com strategy pattern
- ✅ Configuração via variáveis de ambiente ou arquivo .env
- ✅ Tempo de bloqueio configurável
- ✅ Docker e Docker Compose prontos para uso
- ✅ Testes automatizados completos

## 📋 Requisitos

- Docker e Docker Compose
- Go 1.21+ (para desenvolvimento local)

## 🚀 Como Executar

### Usando Docker Compose (Recomendado)

1. Clone o repositório
2. Execute o comando:

```bash
docker-compose up --build
```

O servidor estará disponível em `http://localhost:8080`

### Executando Localmente

1. Certifique-se de que o Redis está rodando
2. Configure as variáveis de ambiente (copie `.env.example` para `.env`)
3. Execute:

```bash
go mod download
go run cmd/server/main.go
```

## ⚙️ Configuração

O sistema pode ser configurado através de variáveis de ambiente ou arquivo `.env`:

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `REDIS_HOST` | Host do Redis | `localhost` |
| `REDIS_PORT` | Porta do Redis | `6379` |
| `REDIS_PASSWORD` | Senha do Redis | (vazio) |
| `REDIS_DB` | Banco de dados do Redis | `0` |
| `RATE_LIMIT_IP_RPS` | Requisições por segundo por IP | `5` |
| `RATE_LIMIT_IP_BLOCK_TIME` | Tempo de bloqueio em segundos para IP | `300` |
| `RATE_LIMIT_TOKENS` | Configuração de tokens (formato: token:rps:blocktime) | (vazio) |
| `SERVER_PORT` | Porta do servidor | `8080` |

### Exemplo de Configuração de Tokens

```env
RATE_LIMIT_TOKENS=abc123:10:300,xyz789:100:600
```

Formato: `TOKEN:RPS:BLOCK_TIME_SECONDS`
- `abc123`: pode fazer 10 requisições por segundo, bloqueado por 300 segundos se exceder
- `xyz789`: pode fazer 100 requisições por segundo, bloqueado por 600 segundos se exceder

## 🔧 Como Funciona

### Fluxo de Requisição

1. **Extração de Identificador**: O middleware extrai o IP (de `X-Forwarded-For`, `X-Real-IP` ou `RemoteAddr`) e o token (header `API_KEY`)

2. **Verificação de Limite**: 
   - Se um token válido for fornecido, usa o limite do token
   - Caso contrário, usa o limite do IP
   - Tokens têm prioridade sobre IPs

3. **Contagem e Validação**:
   - Incrementa contador no Redis com chave única por IP/token
   - Contador expira em 1 segundo (janela deslizante)
   - Se exceder o limite, bloqueia por tempo configurado

4. **Resposta**:
   - ✅ Permitido: Status 200, continua para handler
   - ❌ Bloqueado: Status 429 com mensagem de erro

### Exemplo de Uso

#### Requisição Normal (sem token)

```bash
curl http://localhost:8080/
```

Limite: 5 requisições por segundo (configuração padrão de IP)

#### Requisição com Token

```bash
curl -H "API_KEY: abc123" http://localhost:8080/
```

Limite: 10 requisições por segundo (configuração do token abc123)

#### Excedendo o Limite

Resposta HTTP 429:
```json
{
  "error": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

## 🧪 Testes

Execute os testes com:

```bash
go test ./... -v
```

Execute testes com cobertura:

```bash
go test ./... -cover
```

### Cobertura de Testes

- ✅ Testes unitários para rate limiter
- ✅ Testes de middleware HTTP
- ✅ Testes de configuração
- ✅ Mock de storage para testes isolados
- ✅ Testes de diferentes cenários (IP, token, bloqueio)

## 📁 Estrutura do Projeto

```
.
├── cmd/
│   └── server/
│       └── main.go              # Entry point da aplicação
├── internal/
│   ├── config/
│   │   ├── config.go            # Gerenciamento de configuração
│   │   └── config_test.go       # Testes de configuração
│   ├── limiter/
│   │   ├── limiter.go           # Lógica do rate limiter
│   │   └── limiter_test.go      # Testes do rate limiter
│   ├── middleware/
│   │   ├── ratelimiter.go       # Middleware HTTP
│   │   └── ratelimiter_test.go  # Testes do middleware
│   └── storage/
│       ├── storage.go           # Interface de storage (Strategy Pattern)
│       └── redis.go             # Implementação Redis
├── docker-compose.yml           # Configuração Docker Compose
├── Dockerfile                   # Build da aplicação
├── .env.example                 # Exemplo de configuração
├── go.mod                       # Dependências Go
└── README.md                    # Esta documentação
```

## 🎨 Arquitetura

### Strategy Pattern

O sistema utiliza o padrão Strategy para abstração de storage, permitindo fácil substituição do Redis por outro mecanismo:

```go
type Storage interface {
    Increment(ctx context.Context, key string) (int64, error)
    Get(ctx context.Context, key string) (int64, error)
    SetExpiration(ctx context.Context, key string, expiration time.Duration) error
    IsBlocked(ctx context.Context, key string) (bool, error)
    Block(ctx context.Context, key string, duration time.Duration) error
    Reset(ctx context.Context, key string) error
    Close() error
}
```

Para adicionar um novo storage (ex: Memcached, PostgreSQL):
1. Crie nova struct implementando a interface `Storage`
2. Injete no construtor do `RateLimiter`
3. Pronto! Sem modificar a lógica do limiter

### Separação de Responsabilidades

- **Config**: Gerencia configurações e variáveis de ambiente
- **Storage**: Abstração de persistência (Redis)
- **Limiter**: Lógica de rate limiting pura
- **Middleware**: Integração HTTP, extração de IP/token

## 🧪 Testando com Carga

### Teste Manual de Limite por IP

```bash
# Enviar 10 requisições rapidamente
for i in {1..10}; do
  curl http://localhost:8080/ && echo ""
done
```

As primeiras 5 devem ser bem-sucedidas, a partir da 6ª receberá erro 429.

### Teste com Token

```bash
# Com token abc123 (limite: 10 req/s)
for i in {1..15}; do
  curl -H "API_KEY: abc123" http://localhost:8080/ && echo ""
done
```

As primeiras 10 devem ser bem-sucedidas, a partir da 11ª receberá erro 429.

### Teste com Apache Bench

```bash
# 100 requisições com concorrência de 10
ab -n 100 -c 10 http://localhost:8080/
```

### Teste com wrk

```bash
# Teste de carga por 10 segundos com 2 threads e 10 conexões
wrk -t2 -c10 -d10s http://localhost:8080/
```

## 🔍 Endpoints Disponíveis

- `GET /` - Endpoint de teste que retorna `{"message": "Request successful"}`
- `GET /health` - Health check que retorna `{"status": "healthy"}`

Ambos endpoints estão protegidos pelo rate limiter.

## 📊 Monitoramento

O sistema adiciona headers de resposta para monitoramento:

- `X-RateLimit-Remaining`: Número de requisições restantes na janela atual

## 🛠️ Desenvolvimento

### Adicionar Novo Endpoint

```go
mux.HandleFunc("/seu-endpoint", func(w http.ResponseWriter, r *http.Request) {
    // Seu código aqui
})
```

O middleware já estará aplicado automaticamente.

### Modificar Limites em Tempo de Execução

Edite o arquivo `.env` ou as variáveis de ambiente no `docker-compose.yml` e reinicie:

```bash
docker-compose restart app
```

## 🐛 Troubleshooting

### Redis não conecta

Verifique se o Redis está rodando:
```bash
docker-compose ps redis
```

Veja os logs:
```bash
docker-compose logs redis
```

### Testes falhando

Limpe o cache do Go e execute novamente:
```bash
go clean -testcache
go test ./... -v
```

## 📝 Notas Importantes

1. **Janela de Tempo**: O sistema usa janela deslizante de 1 segundo
2. **Persistência**: Dados são armazenados no Redis com TTL automático
3. **Bloqueio**: Quando bloqueado, o usuário deve aguardar o tempo completo de bloqueio
4. **Prioridade**: Limites de token sempre têm prioridade sobre limites de IP
5. **IP Real**: Sistema detecta IP real através de headers `X-Forwarded-For` e `X-Real-IP`

## 📄 Licença

Este projeto foi desenvolvido como parte de um desafio técnico.

## 👨‍💻 Contato

Para dúvidas ou sugestões, abra uma issue no repositório.
