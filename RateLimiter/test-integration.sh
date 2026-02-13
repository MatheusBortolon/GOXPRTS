#!/bin/bash

# Script de teste de integração do Rate Limiter
# Este script testa diferentes cenários de rate limiting

echo "==================================="
echo "Rate Limiter - Testes de Integração"
echo "==================================="
echo ""

BASE_URL="http://localhost:8080"

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "1. Testando endpoint de health..."
response=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/health)
if [ $response -eq 200 ]; then
    echo -e "${GREEN}✓ Health check passou${NC}"
else
    echo -e "${RED}✗ Health check falhou (Status: $response)${NC}"
    exit 1
fi
echo ""

echo "2. Testando limite por IP (5 req/s)..."
success_count=0
fail_count=0

for i in {1..10}; do
    response=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/)
    if [ $response -eq 200 ]; then
        ((success_count++))
        echo -e "${GREEN}Requisição $i: OK (200)${NC}"
    elif [ $response -eq 429 ]; then
        ((fail_count++))
        echo -e "${RED}Requisição $i: Bloqueada (429)${NC}"
    fi
    sleep 0.1
done

echo ""
echo "Resultado: $success_count sucesso, $fail_count bloqueadas"
if [ $success_count -eq 5 ] && [ $fail_count -eq 5 ]; then
    echo -e "${GREEN}✓ Limite por IP funcionando corretamente${NC}"
else
    echo -e "${YELLOW}⚠ Resultado inesperado. Esperado: 5 sucesso, 5 bloqueadas${NC}"
fi
echo ""

# Aguardar reset do contador (1 segundo)
echo "Aguardando reset do contador..."
sleep 2

echo "3. Testando limite por token abc123 (10 req/s)..."
success_count=0
fail_count=0

for i in {1..15}; do
    response=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: abc123" $BASE_URL/)
    if [ $response -eq 200 ]; then
        ((success_count++))
        echo -e "${GREEN}Requisição $i: OK (200)${NC}"
    elif [ $response -eq 429 ]; then
        ((fail_count++))
        echo -e "${RED}Requisição $i: Bloqueada (429)${NC}"
    fi
    sleep 0.1
done

echo ""
echo "Resultado: $success_count sucesso, $fail_count bloqueadas"
if [ $success_count -eq 10 ] && [ $fail_count -eq 5 ]; then
    echo -e "${GREEN}✓ Limite por token funcionando corretamente${NC}"
else
    echo -e "${YELLOW}⚠ Resultado inesperado. Esperado: 10 sucesso, 5 bloqueadas${NC}"
fi
echo ""

# Aguardar reset do contador
echo "Aguardando reset do contador..."
sleep 2

echo "4. Testando prioridade de token sobre IP..."
echo "Fazendo 8 requisições com token (limite 10) - deve passar"
success_count=0

for i in {1..8}; do
    response=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: xyz789" $BASE_URL/)
    if [ $response -eq 200 ]; then
        ((success_count++))
    fi
    sleep 0.1
done

if [ $success_count -eq 8 ]; then
    echo -e "${GREEN}✓ Token xyz789 permite mais requisições que limite de IP${NC}"
else
    echo -e "${RED}✗ Falha no teste de prioridade${NC}"
fi
echo ""

echo "5. Testando resposta de bloqueio..."
# Aguardar reset
sleep 2

# Esgotar o limite
for i in {1..6}; do
    curl -s -o /dev/null $BASE_URL/
    sleep 0.1
done

response=$(curl -s $BASE_URL/)
if echo "$response" | grep -q "you have reached the maximum number of requests"; then
    echo -e "${GREEN}✓ Mensagem de bloqueio correta${NC}"
else
    echo -e "${RED}✗ Mensagem de bloqueio incorreta${NC}"
    echo "Resposta: $response"
fi
echo ""

echo "==================================="
echo "Testes concluídos!"
echo "==================================="
