package messagebroker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/zap"
)

const (
	KAFKA_BROKER_URL_ENV_NAME      = "KAFKA_BROKER_URL"
	KAFKA_CLIENT_USERNAME_ENV_NAME = "KAFKA_CLIENT_USERNAME"
	KAFKA_CLIENT_PASSWORD_ENV_NAME = "KAFKA_CLIENT_PASSWORD"
	KAFKA_MAX_MESSAGE_BYTES_KEY    = "KAFKA_MAX_MESSAGE_SIZE_BYTES"
	KAFKA_BROKER_GROUP_ID_KEY      = "KAFKA_BROKER_GROUP_ID"
)

type KafkaCredentials struct {
	Address  string
	Username string
	Password string
}

func NewKafkaCredentials() *KafkaCredentials {
	envs.ValidateEnvs("%s not set", []string{
		KAFKA_BROKER_URL_ENV_NAME,
		KAFKA_CLIENT_USERNAME_ENV_NAME,
		KAFKA_CLIENT_PASSWORD_ENV_NAME,
	})

	return &KafkaCredentials{
		Address:  os.Getenv(KAFKA_BROKER_URL_ENV_NAME),
		Username: os.Getenv(KAFKA_CLIENT_USERNAME_ENV_NAME),
		Password: os.Getenv(KAFKA_CLIENT_PASSWORD_ENV_NAME),
	}
}

type KafkaMessageBroker struct {
	logger *zap.Logger
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaMessageBroker(addr string, topic string, username string, password string, logger *zap.Logger) *KafkaMessageBroker {

	envs.ValidateEnvs("No max message size set for kafka broker", []string{KAFKA_MAX_MESSAGE_BYTES_KEY})
	envs.ValidateEnvs("No consumer group id set for kafka broker", []string{KAFKA_BROKER_GROUP_ID_KEY})

	kafkaMaxMessageBytes := os.Getenv(KAFKA_MAX_MESSAGE_BYTES_KEY)
	kafkaMaxMessageBytesInt, err := strconv.Atoi(kafkaMaxMessageBytes)

	kafkaBrokerGroupId := os.Getenv(KAFKA_BROKER_GROUP_ID_KEY)

	if err != nil {
		panic("Kafka max message size is not a number")
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: username, Password: password}},
		BatchBytes:             int64(kafkaMaxMessageBytesInt),
		// BatchBytes: 0,
		// BatchSize:  0,
	}

	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		panic("Failed to set sasl mechanism for logs ingestion kafka queue")
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
	}

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        []string{addr},
			Topic:          topic,
			MaxBytes:       10e8,
			GroupID:        kafkaBrokerGroupId,
			Dialer:         dialer,
			CommitInterval: time.Second,
			QueueCapacity:  0,
			MinBytes:       0,
			MaxWait:        time.Duration(time.Second),
			// CommitInterval: ,
		},
	)

	return &KafkaMessageBroker{
		writer: writer,
		reader: reader,
		logger: logger,
	}
}

func (b *KafkaMessageBroker) Publish(ctx context.Context, key []byte, value []byte) error {
	err := b.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})

	if err != nil {
		b.logger.Error("Failed to publish a message", zap.Error(err),
			zap.String("key", string(key)),
		)
		return err
	}

	return nil
}

func (b *KafkaMessageBroker) Subscribe(ctx context.Context, messages chan<- []byte, errors chan<- error) {
	for {
		msg, err := b.reader.ReadMessage(ctx)
		fmt.Println("read message: ", msg)

		if err != nil {
			b.logger.Error("Failed to read message", zap.Error(err))
			errors <- err
			continue
		}

		messages <- msg.Value
	}
}
func (b *KafkaMessageBroker) CloseReader() error {
	return b.reader.Close()
}

var _ MessageBroker[any] = &KafkaJsonMessageBroker[any]{}
