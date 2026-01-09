# Stress Test CLI em Go

Sistema CLI para realizar testes de carga em um serviço web.

## Funcionalidades

- ✅ Testes de carga HTTP parametrizáveis
- ✅ Controle de concorrência
- ✅ Relatório detalhado com distribuição de status codes
- ✅ Métricas de performance (requisições por segundo)
- ✅ Containerização com Docker

## Instalação

### Pré-requisitos
- Go 1.21+
- Docker (opcional)

### Build Local

```bash
go build -o stresstest
```

### Build Docker

```bash
docker build -t stresstest:latest .
```

## Uso

### Localmente

```bash
./stresstest --url=http://google.com --requests=1000 --concurrency=10
```

### Via Docker

```bash
docker run stresstest:latest --url=http://google.com --requests=1000 --concurrency=10
```

## Parâmetros

- `--url` (obrigatório): URL do serviço a ser testado
- `--requests` (obrigatório): Número total de requests a realizar
- `--concurrency` (obrigatório): Número de chamadas simultâneas (goroutines)

## Exemplo de Saída

```
============================================================
RELATÓRIO DE TESTE DE CARGA
============================================================
Tempo total: 2.45s
Total de requests: 1000
Requests por segundo: 408.16
Erros: 0

Distribuição de códigos HTTP:
------------------------------------------------------------
  Status 200: 1000
============================================================
```

## Arquitetura

- **main.go**: Contém toda a lógica do aplicativo
  - `parseFlags()`: Parse dos argumentos CLI
  - `validateConfig()`: Validação dos parâmetros
  - `runStressTest()`: Orquestração dos testes
  - `worker()`: Função executada por cada goroutine
  - `printReport()`: Formatação e exibição do relatório

## Detalhes da Implementação

### Concorrência
- Utiliza goroutines do Go para executar requests em paralelo
- Canal (channel) para distribuição de work entre workers
- WaitGroup para sincronização de conclusão

### Medições
- Tempo total: Diferença entre início e fim dos testes
- Taxa de requisições: Total de requisições / tempo em segundos
- Status codes: Mapeamento de cada código HTTP recebido
- Erros de conexão: Rastreamento de falhas de rede/timeout

### Performance
- HTTP Client com timeout configurado (30s)
- Reutilização de conexões via HTTP client
- Goroutines leves para melhor escalabilidade
