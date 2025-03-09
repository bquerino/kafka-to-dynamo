package main

import (
	"log"

	"github.com/bquerino/kafka-to-dynamo/internal/config"
	"github.com/bquerino/kafka-to-dynamo/internal/consumer"
	"github.com/bquerino/kafka-to-dynamo/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Cria o repositório para DynamoDB
	dynamoRepo := repository.NewDynamoRepo(cfg)

	// Cria o consumidor Kafka
	kafkaConsumer, err := consumer.NewKafkaConsumer(cfg, dynamoRepo)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	// Inicia o loop de consumo
	kafkaConsumer.Run()
}
