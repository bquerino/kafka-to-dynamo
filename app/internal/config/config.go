// internal/config/config.go
package config

import "os"

// Config contém as configurações da aplicação.
type Config struct {
	KafkaBrokers      string
	KafkaTopic        string
	SchemaRegistryURL string // Pode ser usado para outras integrações, se necessário
	DynamoEndpoint    string
	DynamoTableName   string
}

// Load carrega as configurações via variáveis de ambiente.
func Load() (*Config, error) {
	return &Config{
		KafkaBrokers:      os.Getenv("KAFKA_BROKERS"),       // Ex.: "localhost:9092"
		KafkaTopic:        os.Getenv("KAFKA_TOPIC"),         // Ex.: "payments-topic"
		SchemaRegistryURL: os.Getenv("SCHEMA_REGISTRY_URL"), // Ex.: "http://localhost:8081"
		DynamoEndpoint:    os.Getenv("DYNAMO_ENDPOINT"),     // Ex.: "http://localhost:4566"
		DynamoTableName:   os.Getenv("DYNAMO_TABLE_NAME"),   // Ex.: "Pagamentos"
	}, nil
}
