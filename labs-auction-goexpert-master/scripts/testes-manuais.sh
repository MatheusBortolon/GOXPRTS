#!/usr/bin/env bash
set -euo pipefail

# Roteiro de testes manuais para validacao do fechamento automatico.
# Requisitos: API em http://localhost:8080 e AUCTION_INTERVAL curto (ex: 5s).

USER_ID="00000000-0000-4000-8000-000000000001"

echo "1) Criar leilao"
AUCTION_CREATE_RESPONSE=$(curl -s -i -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{"product_name":"TV 4K TESTE","category":"Eletronicos","description":"TV nova com garantia","condition":1}')

echo "$AUCTION_CREATE_RESPONSE" | head -n 1

echo "2) Buscar leilao criado (status=0)"
AUCTIONS_JSON=$(curl -s "http://localhost:8080/auction?status=0")

echo "$AUCTIONS_JSON"

AUCTION_ID=$(echo "$AUCTIONS_JSON" | tr '{' '\n' | grep '"product_name":"TV 4K TESTE"' | sed -n 's/.*"id":"\([^"]*\)".*/\1/p' | head -n 1)
if [ -z "$AUCTION_ID" ]; then
  echo "Erro: nao foi possivel encontrar o leilao criado."
  exit 1
fi

echo "Leilao ID: $AUCTION_ID"

echo "3) Confirmar ativo (status=0)"
curl -s "http://localhost:8080/auction/$AUCTION_ID"

echo "4) Enviar lances em intervalos aleatorios (antes do fechamento)"
for i in 1 2 3 4 5; do
  sleep_time=$(( (RANDOM % 3) + 1 ))
  amount=$(( (RANDOM % 200) + 50 ))
  echo "- Lance $i: aguardando ${sleep_time}s, valor ${amount}.00"
  sleep "$sleep_time"
  curl -s -i -X POST http://localhost:8080/bid \
    -H "Content-Type: application/json" \
    -d '{"user_id":"'"$USER_ID"'","auction_id":"'"$AUCTION_ID"'","amount":'"$amount"'.0}' | head -n 1
done

echo "5) Aguardar fechamento automatico"
sleep 7

echo "6) Confirmar fechado (status=1)"
curl -s "http://localhost:8080/auction/$AUCTION_ID"

echo "7) Tentar criar lance apos fechar (nao deve persistir)"
curl -s -i -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'","auction_id":"'"$AUCTION_ID"'","amount":150.0}' | head -n 1

echo "8) Conferir lances (esperado: apenas os antes do fechamento)"
curl -s "http://localhost:8080/bid/$AUCTION_ID"