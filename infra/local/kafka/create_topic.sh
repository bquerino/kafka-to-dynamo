#!/bin/bash

# Configurações
KAFKA_REST_URL="http://localhost:8082"  # URL do Kafka REST Proxy
TOPIC_NAME="payments-done"
PARTITIONS=3  # Número de partições
REPLICATION_FACTOR=1  # Fator de replicação

# Endpoint para listar clusters
CLUSTERS_URL="$KAFKA_REST_URL/v3/clusters"

# Obtém o ID do cluster
CLUSTER_ID=$(curl -s "$CLUSTERS_URL" | jq -r '.data[0].cluster_id')

if [ -z "$CLUSTER_ID" ]; then
  echo "Erro: Não foi possível obter o ID do cluster."
  exit 1
fi

echo "ID do cluster encontrado: $CLUSTER_ID"

# Endpoint para criar tópicos
CREATE_TOPIC_URL="$KAFKA_REST_URL/v3/clusters/$CLUSTER_ID/topics"

# Corpo da requisição (JSON)
PAYLOAD=$(cat <<EOF
{
  "topic_name": "$TOPIC_NAME",
  "partitions_count": $PARTITIONS,
  "replication_factor": $REPLICATION_FACTOR,
  "configs": [
    {"name": "cleanup.policy", "value": "delete"},
    {"name": "retention.ms", "value": "604800000"}
  ]
}
EOF
)

# Headers
HEADERS=(
  "Content-Type: application/json"
)

# Faz a requisição POST para criar o tópico
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$CREATE_TOPIC_URL" \
  -H "${HEADERS[0]}" \
  -d "$PAYLOAD")

# Verifica a resposta
if [ "$RESPONSE" -eq 201 ]; then
  echo "Tópico '$TOPIC_NAME' criado com sucesso!"
else
  echo "Erro ao criar tópico. Código de resposta: $RESPONSE"
  # Exibe a resposta completa em caso de erro
  curl -X POST "$CREATE_TOPIC_URL" -H "${HEADERS[0]}" -d "$PAYLOAD"
fi