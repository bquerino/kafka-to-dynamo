package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"

	"github.com/bquerino/kafka-to-dynamo/internal/config"
	"github.com/bquerino/kafka-to-dynamo/internal/domain"
	"github.com/bquerino/kafka-to-dynamo/internal/repository"
)

const avroSchema = `{
  "type": "record",
  "name": "PaymentRecord",
  "namespace": "com.example.payments",
  "fields": [
    {"name": "PaymentID", "type": "string"},
    {"name": "CustomerID", "type": "string"},
    {"name": "PaymentTimestamp", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "TransactionValue", "type": "double"}
  ]
}`

// KafkaConsumer encapsula o leitor do Kafka, o repositório e o codec Avro.
type KafkaConsumer struct {
	reader *kafka.Reader
	repo   repository.DynamoRepo
	codec  *goavro.Codec
}

// NewKafkaConsumer cria e configura o consumer Kafka.
func NewKafkaConsumer(cfg *config.Config, repo repository.DynamoRepo) (*KafkaConsumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBrokers}, // Ex.: "localhost:9092"
		Topic:    cfg.KafkaTopic,             // Ex.: "payments-topic"
		GroupID:  "group-payment",
		MinBytes: 10e3, // 10 KB
		MaxBytes: 10e6, // 10 MB
	})

	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		return nil, fmt.Errorf("error creating Avro codec: %w", err)
	}

	return &KafkaConsumer{
		reader: r,
		repo:   repo,
		codec:  codec,
	}, nil
}

// Run inicia o loop de consumo, decodifica a mensagem Avro e persiste o pagamento.
func (kc *KafkaConsumer) Run() {
	ctx := context.Background()
	log.Println("Starting Kafka consumer...")

	for {
		msg, err := kc.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Error fetching message from Kafka: %v", err)
			break
		}

		// Se o payload estiver com o header do Schema Registry (5 bytes: magic byte + 4-byte schema ID), remova-os
		var payload []byte
		if len(msg.Value) > 5 && msg.Value[0] == 0 {
			payload = msg.Value[5:]
		} else {
			payload = msg.Value
		}

		// Decodifica o payload Avro
		native, _, err := kc.codec.NativeFromBinary(payload)
		if err != nil {
			log.Printf("Error decoding Avro message: %v", err)
			_ = kc.reader.CommitMessages(ctx, msg)
			continue
		}

		// Converte para map[string]interface{}
		recordMap, ok := native.(map[string]interface{})
		if !ok {
			log.Printf("Error asserting Avro record to map")
			_ = kc.reader.CommitMessages(ctx, msg)
			continue
		}

		// Extrai os campos
		paymentID, _ := recordMap["PaymentID"].(string)
		customerID, _ := recordMap["CustomerID"].(string)

		var paymentTimestamp int64
		switch v := recordMap["PaymentTimestamp"].(type) {
		case int64:
			paymentTimestamp = v
		case int32:
			paymentTimestamp = int64(v)
		case float64:
			paymentTimestamp = int64(v)
		case time.Time:
			paymentTimestamp = v.UnixNano() / int64(1e6)
		default:
			log.Printf("Unexpected type for PaymentTimestamp: %T", v)
		}

		var transactionValue float64
		switch v := recordMap["TransactionValue"].(type) {
		case float64:
			transactionValue = v
		case float32:
			transactionValue = float64(v)
		default:
			log.Printf("Unexpected type for TransactionValue: %T", v)
		}

		// Cria o objeto Payment conforme o domínio
		payment := &domain.Payment{
			PaymentID:        paymentID,
			CustomerID:       customerID,
			PaymentTimestamp: paymentTimestamp,
			TransactionValue: transactionValue,
		}

		// Persiste o pagamento no DynamoDB
		if err := kc.repo.SavePayment(ctx, payment); err != nil {
			log.Printf("Error saving payment to DynamoDB: %v", err)
		} else {
			log.Printf("Payment saved: %+v", payment)
		}

		// Comita a mensagem para confirmar o processamento
		if err := kc.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}

// Close encerra o leitor Kafka.
func (kc *KafkaConsumer) Close() error {
	log.Println("Closing Kafka consumer.")
	return kc.reader.Close()
}
