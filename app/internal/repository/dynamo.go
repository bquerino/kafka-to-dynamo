package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bquerino/kafka-to-dynamo/internal/config"
	"github.com/bquerino/kafka-to-dynamo/internal/domain"
)

// DynamoRepo define as operações de persistência.
type DynamoRepo interface {
	SavePayment(ctx context.Context, payment *domain.Payment) error
}

type dynamoRepo struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoRepo cria a instância do repositório para DynamoDB.
func NewDynamoRepo(cfg *config.Config) DynamoRepo {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           cfg.DynamoEndpoint, // Ex.: "http://localhost:4566"
					SigningRegion: "us-east-1",
				}, nil
			}),
		),
	)
	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(awsCfg)
	return &dynamoRepo{
		client:    client,
		tableName: cfg.DynamoTableName,
	}
}

// SavePayment persiste um pagamento no DynamoDB.
func (r *dynamoRepo) SavePayment(ctx context.Context, p *domain.Payment) error {
	// Monta a chave primária incluindo os identificadores
	pk := fmt.Sprintf("CUSTOMERID#%s#PAYMENTID#%s", p.CustomerID, p.PaymentID)
	// A chave de ordenação incorpora o timestamp
	sk := fmt.Sprintf("Timestamp#%d", p.PaymentTimestamp)

	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item: map[string]types.AttributeValue{
			"CustomerID#PaymentID": &types.AttributeValueMemberS{Value: pk},
			"PaymentEventDate":     &types.AttributeValueMemberS{Value: sk},
			"TransactionValue":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", p.TransactionValue)},
		},
	})
	return err
}
