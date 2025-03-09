#!/bin/bash

# Configurações
KAFKA_REST_URL="http://localhost:8082"  # URL do Kafka REST Proxy
SCHEMA_REGISTRY_URL="http://localhost:8081"  # URL do Schema Registry
TOPIC_NAME="payments-done"
SUBJECT="${TOPIC_NAME}-value"  # Assunto no Schema Registry

# Obtém o ID do schema Avro mais recente
SCHEMA_ID=$(curl -s "${SCHEMA_REGISTRY_URL}/subjects/${SUBJECT}/versions/latest" | jq -r '.id')

if [ -z "$SCHEMA_ID" ]; then
  echo "Erro: Não foi possível obter o ID do schema Avro."
  exit 1
fi

echo "ID do schema Avro: $SCHEMA_ID"

# Endpoint para produzir mensagens
PRODUCE_URL="${KAFKA_REST_URL}/topics/${TOPIC_NAME}"

# Dados da mensagem (exemplo)
PAYLOAD_DATA=$(cat <<EOF
{
  "value_schema_id": $SCHEMA_ID,
  "records": [
    {
      "value": {
        "PaymentID": "pay_$(date +%s)",
        "CustomerID": "cust_$(shuf -i 1000-9999 -n 1)",
        "PaymentTimestamp": $(date +%s%3N),
        "TransactionValue": $(awk -v min=10 -v max=1000 'BEGIN{srand(); printf "%.2f", min+rand()*(max-min)}')
      }
    }
  ]
}
EOF
)

# Headers
HEADERS=(
  "Content-Type: application/vnd.kafka.avro.v2+json"
)

# Faz a requisição POST para produzir a mensagem
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$PRODUCE_URL" \
  -H "${HEADERS[0]}" \
  -d "$PAYLOAD_DATA")

# Verifica a resposta
if [ "$RESPONSE" -eq 200 ]; then
  echo "Mensagem produzida com sucesso no tópico '$TOPIC_NAME'!"
else
  echo "Erro ao produzir mensagem. Código de resposta: $RESPONSE"
  curl -X POST "$PRODUCE_URL" -H "${HEADERS[0]}" -d "$PAYLOAD_DATA"
fi