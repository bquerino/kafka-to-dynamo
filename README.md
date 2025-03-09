# **Kafka -> Dynamo**
An Go Lang app that get an message from Kafka and push to AWS DynamoDB.

## **Setup Local**

Esse projeto faz o uso do LocalStack, para isso certifique-se que tenha o AWS CLI previamente instalado.

### **Inicie o `docker-compose`**

Acesse a pasta `infra/local/` e rode o comando:

```bash
docker-compose up -d
```

> O parâmetro -d roda os containers em modo detached - em background para que os logs não apareçam após a execução. Caso prefira visualizar, remova o parâmetro.

### **Crie seu perfil `localstack`**

```bash
aws configure --profile localstack
```

```bash
AWS Access Key ID: test (ou qualquer valor fictício)
AWS Secret Access Key: test (ou qualquer valor fictício)
Default region name: us-east-1 (ou outra região de sua preferência)
Default output format: json (ou o formato que quiser)
```

> Isso criará (ou atualizará) a entrada `[localstack]` dentro de `~/.aws/credentials` e `[profile localstack]` em `~/.aws/config`.

Feito isso utilize sempre o parâmetro `--endpoint-url` para apontar para o LocalStack, todo comando inclua: `--endpoint-url http://localhost:4566`. Por exemplo:

```bash
aws --profile localstack --endpoint-url http://localhost:4566 dynamodb list-tables
```

### **Crie os recursos AWS**

Acesse a pasta `infra/local/terraform` e rode o script:

```bash
./run_tf.sh
```

### **Crie os componentes do Kafka**

Acesse a pasta `infra/local/kafka` e rode os scripts:

```bash
./create_topic.sh
./schema_registry.sh
```

## **Rodando a aplicação**

Acesse a pasta `/app/cmd/kafka_consumer/` e rode o script:

```bash
./start.sh
```

Após isso envie uma mensagem usando o script da pasta `/infra/local/kafka`

```bash
./send_message.sh
```

Verifique se a mensagem foi gravada no DynamoDB usando a seguinte consulta:

```bash
aws dynamodb scan --table-name payments --endpoint-url http://localhost:4566 --region us-east-1
```

Para verificar se a mensagem chegou ao Kafka acesse o Kafka ControlCenter através da URL: http://localhost:9021/

---
## **References**

- [Confluent Platform](https://docs.confluent.io/platform/current/get-started/platform-quickstart.html)
- [LocalStack](https://www.localstack.cloud/)