# Real-Time Payment Processing: Kafka to DynamoDB Pipeline 🔄💸

*A high-throughput pipeline for processing payment events in real-time*

## Features ✨
- **Real-Time Event Consumption**: Process Kafka messages with low latency
- **Schema Enforcement**: Avro schema validation via Schema Registry
- **DynamoDB Integration**: Store payment records with automatic retries
- **Local Development**: Full local stack with LocalStack + Confluent
- **Infra-as-Code**: Terraform-provisioned AWS resources

---

## Tech Stack 🛠️
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Apache Kafka](https://img.shields.io/badge/Apache_Kafka-3.6+-231F20?logo=apache-kafka)
![Terraform](https://img.shields.io/badge/Terraform-1.5+-7B42BC?logo=terraform)
![LocalStack](https://img.shields.io/badge/LocalStack-3.0+-5A0FC8?logo=localstack)

---

## Prerequisites 📋
- [Go 1.21+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [AWS CLI v2](https://aws.amazon.com/cli/)
- [Terraform 1.5+](https://www.terraform.io/downloads)
- [jq](https://stedolan.github.io/jq/download/) (JSON processor)

---

## Local Development Setup 🛠️

### 1. Start Local Infrastructure
```bash
cd infra/local/
docker-compose up -d --build
```

## Services Overview 🖥️

| Service              | Port  | URL                          |
|----------------------|-------|------------------------------|
| Kafka Broker         | 9092  | `PLAINTEXT://localhost:9092` |
| Schema Registry      | 8081  | `http://localhost:8081`      |
| Kafka Connect        | 8083  | `http://localhost:8083`      |
| LocalStack (Dynamo)  | 4566  | `http://localhost:4566`      |
| Control Center       | 9021  | `http://localhost:9021`      |

### 2. Configure AWS CLI Profile
```bash
aws configure --profile localstack
```

Use these dummy credentials:

```text
AWS Access Key ID: test
AWS Secret Access Key: test
Default region: us-east-1
Output format: json
```

### 3. Provision AWS Resources

```bash
cd infra/local/terraform/
./run_tf.sh
```

**Created Resources:**

* DynamoDB Table: payments

### 4. Initialize Kafka Environment

```bash
cd infra/local/kafka/
./create_topic.sh          # Creates 'payments-done' topic
./schema_registry.sh       # Registers Avro schema
```
--- 
## Running the Application ▶️

### Start Kafka Consumer
```bash
cd app/cmd/kafka_consumer/
./start.sh
```

### Environment Variables:

```env
KAFKA_BROKERS=localhost:9092
SCHEMA_REGISTRY_URL=http://localhost:8081
DYNAMO_TABLE=payments
AWS_ENDPOINT=http://localhost:4566
```

### Produce Test Messages
```bash
cd infra/local/kafka/
./send_message.sh
```

### Sample Message:

```json
{
  "PaymentID": "pay_1691439325",
  "CustomerID": "cust_8823",
  "PaymentTimestamp": 1691439325123,
  "TransactionValue": 249.99
}
```
---

## Verification ✔️
### Check DynamoDB Records
```bash
aws dynamodb scan \
  --table-name payments \
  --endpoint-url http://localhost:4566 \
  --region us-east-1 \
  --profile localstack
```

### Monitor Kafka Metrics

Access Control Center: http://localhost:9021

* View consumer lag
* Inspect message schemas
* Monitor throughput metrics

---

## Project Structure 📁
```text
Copy
├── app/
│   └── cmd/
│       └── kafka_consumer/   # Main application
|   └── internal/
        └── config/     # App configs
        └── consumer/   # Kafka consumer implementation
        └── domain/     # Anemic model of Payments
        └── repository/ # Persistence logic
├── infra/
│   ├── local/
│   │   ├── kafka/           # Topic/schema scripts
│   │   ├── terraform/       # LocalStack resources
│   │   └── docker-compose.yml
```

---

## Troubleshooting 🚨

### Common Issues:

* Connection refused errors: Verify Docker services are running
* Schema registry 404: Run schema_registry.sh after topic creation

### Logs Inspection:

```bash
docker-compose logs -f kafka dynamodb
```
---

## References 📚

- [Confluent Kafka Docs](https://docs.confluent.io/platform/current/get-started/platform-quickstart.html)
- [LocalStack AWS Coverage](https://www.localstack.cloud/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Go AWS SDK v2](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/welcome.html)
- [Kafka-Go Library](https://github.com/segmentio/kafka-go)