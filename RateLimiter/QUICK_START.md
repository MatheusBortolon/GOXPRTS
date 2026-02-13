# 🚀 Guia Rápido de Início

## Iniciar o Sistema

```bash
docker-compose up --build
```

## Testar o Sistema

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Requisição Normal
```bash
curl http://localhost:8080/
```

### 3. Requisição com Token
```bash
curl -H "API_KEY: abc123" http://localhost:8080/
```

### 4. Testar Limite (5 requisições rápidas)
```bash
for i in {1..10}; do curl http://localhost:8080/ && echo ""; done
```

Resultado esperado: Primeiras 5 passam, restantes retornam 429.

### 5. Testar com Token (10 requisições)
```bash
for i in {1..15}; do curl -H "API_KEY: abc123" http://localhost:8080/ && echo ""; done
```

Resultado esperado: Primeiras 10 passam, restantes retornam 429.

## Executar Testes Automatizados

### Testes Unitários
```bash
go test ./... -v
```

### Testes de Integração (Linux/Mac)
```bash
bash test-integration.sh
```

### Testes de Integração (Windows)
```powershell
.\test-integration.ps1
```

## Configurar Limites Personalizados

Edite o arquivo `.env` ou `docker-compose.yml`:

```env
# Limite de IP: 10 requisições por segundo
RATE_LIMIT_IP_RPS=10

# Tempo de bloqueio: 60 segundos (1 minuto)
RATE_LIMIT_IP_BLOCK_TIME=60

# Tokens personalizados
RATE_LIMIT_TOKENS=mytoken:50:120,premium:200:300
```

Reinicie o container:
```bash
docker-compose restart app
```

## Verificar Logs

```bash
docker-compose logs -f app
```

## Parar o Sistema

```bash
docker-compose down
```

## Limpar Tudo (incluindo volumes)

```bash
docker-compose down -v
```

## Estrutura de Resposta

### Sucesso (200 OK)
```json
{
  "message": "Request successful"
}
```

### Bloqueado (429 Too Many Requests)
```json
{
  "error": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

## Configurações Padrão

- **IP**: 5 requisições/segundo, bloqueio de 300 segundos (5 minutos)
- **Token abc123**: 10 requisições/segundo, bloqueio de 300 segundos
- **Token xyz789**: 100 requisições/segundo, bloqueio de 600 segundos (10 minutos)
- **Porta**: 8080

## Comandos Úteis com Make

```bash
make help              # Mostra todos os comandos
make docker-build      # Builda e inicia containers
make docker-logs       # Mostra logs em tempo real
make test              # Executa testes unitários
make test-cover        # Gera relatório de cobertura
```

## Troubleshooting

### Porta 8080 já está em uso
Mude a porta no `docker-compose.yml`:
```yaml
ports:
  - "9090:8080"  # Usar porta 9090 localmente
```

### Redis não conecta
Verifique se está rodando:
```bash
docker-compose ps redis
docker-compose logs redis
```

### Container não inicia
Veja os logs:
```bash
docker-compose logs app
```

## Próximos Passos

1. ✅ Sistema rodando em http://localhost:8080
2. ✅ Redis funcionando
3. 🔧 Personalize os limites conforme sua necessidade
4. 🧪 Execute os testes de integração
5. 📊 Monitore via logs
