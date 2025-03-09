#!/bin/bash

# Define as variáveis de ambiente necessárias
export KAFKA_BROKERS="localhost:9092"
export KAFKA_TOPIC="payments-done"
export SCHEMA_REGISTRY_URL="http://localhost:8081"
export DYNAMO_ENDPOINT="http://localhost:4566"
export DYNAMO_TABLE_NAME="Payments"

# Inicia a aplicação
go run main.go
