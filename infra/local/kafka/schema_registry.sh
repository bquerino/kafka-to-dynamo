#!/bin/bash

# Configurações
SCHEMA_REGISTRY_URL="http://localhost:8081"  # URL do Schema Registry
TOPIC_NAME="payments-done"
SUBJECT="${TOPIC_NAME}-value"  # Assunto no Schema Registry

# Schema Avro para pagamento (definido como string JSON escapada)
AVRO_SCHEMA=$(cat <<EOF
{
  "type": "record",
  "name": "PaymentRecord",
  "namespace": "com.example.payments",
  "fields": [
    {"name": "PaymentID", "type": "string"},
    {"name": "CustomerID", "type": "string"},
    {"name": "PaymentTimestamp", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "TransactionValue", "type": "double"}
  ]
}
EOF
)

# Endpoint para registrar o schema
REGISTER_SCHEMA_URL="$SCHEMA_REGISTRY_URL/subjects/$SUBJECT/versions"

# Converter o schema para uma string JSON escapada
ESCAPED_AVRO_SCHEMA=$(echo "$AVRO_SCHEMA" | jq -c | jq -R)

# Corpo da requisição (schema como string)
PAYLOAD=$(cat <<EOF
{
  "schema": $ESCAPED_AVRO_SCHEMA
}
EOF
)

# Headers
HEADERS=(
  "Content-Type: application/vnd.schemaregistry.v1+json"
)

# Faz a requisição POST para registrar o schema
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$REGISTER_SCHEMA_URL" \
  -H "${HEADERS[0]}" \
  -d "$PAYLOAD")

# Verifica a resposta
if [ "$RESPONSE" -eq 200 ]; then
  SCHEMA_ID=$(curl -s -X POST "$REGISTER_SCHEMA_URL" -H "${HEADERS[0]}" -d "$PAYLOAD" | jq -r '.id')
  echo "Schema Avro registrado com sucesso!"
  echo "ID do Schema: $SCHEMA_ID"
else
  echo "Erro ao registrar schema. Código de resposta: $RESPONSE"
  curl -X POST "$REGISTER_SCHEMA_URL" -H "${HEADERS[0]}" -d "$PAYLOAD"
fi