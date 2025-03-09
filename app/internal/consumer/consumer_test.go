package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/bquerino/kafka-to-dynamo/internal/domain"
	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

func (m *MockReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SavePayment(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func TestKafkaConsumer_Run(t *testing.T) {
	mockReader := new(MockReader)
	mockRepo := new(MockRepo)

	codec, _ := goavro.NewCodec(avroSchema)
	kc := &KafkaConsumer{
		reader: mockReader,
		repo:   mockRepo,
		codec:  codec,
	}

	paymentData := map[string]interface{}{
		"PaymentID":        "pay_123",
		"CustomerID":       "cust_456",
		"PaymentTimestamp": int64(1710000000000),
		"TransactionValue": 123.45,
	}

	binaryData, _ := codec.BinaryFromNative(nil, paymentData)
	msg := kafka.Message{Value: binaryData}

	mockReader.On("FetchMessage", mock.Anything).Return(msg, nil).Once()
	mockReader.On("FetchMessage", mock.Anything).Return(kafka.Message{}, errors.New("EOF")).Once()
	mockReader.On("CommitMessages", mock.Anything, []kafka.Message{msg}).Return(nil)
	mockRepo.On("SavePayment", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)

	kc.Run()

	mockReader.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestKafkaConsumer_Close(t *testing.T) {
	mockReader := new(MockReader)
	codec, _ := goavro.NewCodec(avroSchema)

	kc := &KafkaConsumer{
		reader: mockReader,
		codec:  codec,
	}

	mockReader.On("Close").Return(nil)

	err := kc.Close()
	assert.NoError(t, err)

	mockReader.AssertExpectations(t)
}
