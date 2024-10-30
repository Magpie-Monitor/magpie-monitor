package brokers

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"time"
)

var JOB_SCHEDULED_TOPIC_KEY = "REPORT_GENERATED_BROKER_TOPIC"

type OpenAiJob struct {
	CompletionRequests []openai.CompletionRequest
}

type JobScheduled[T any] struct {
	Id            string `json:"id" bson:"_id"`
	CorrelationId string `json:"correlationId"`
	ScheduledAtMs int64  `json:"scheduledAtMs"`
	Data          any    `json:"data"`
}

func NewJobScheduled[T any](id string, correlationId string, data T) *JobScheduled[T] {

	return &JobScheduled[T]{
		Id:            id,
		CorrelationId: correlationId,
		ScheduledAtMs: time.Now().UnixMilli(),
		Data:          data,
	}
}

func NewOpenAiJobScheduledBroker(logger *zap.Logger) *messagebroker.KafkaJsonMessageBroker[JobScheduled[OpenAiJob]] {

	envs.ValidateEnvs(
		"address/username/password/topic for JobScheduledBroker is not set",
		[]string{
			MESSAGE_BROKER_ADDRESS_KEY,
			MESSAGE_BROKER_PASSWORD_KEY,
			MESSAGE_BROKER_USERNAME_KEY,
			REPORT_GENERATED_TOPIC_KEY,
		},
	)

	username := os.Getenv(MESSAGE_BROKER_USERNAME_KEY)
	password := os.Getenv(MESSAGE_BROKER_PASSWORD_KEY)
	address := os.Getenv(MESSAGE_BROKER_ADDRESS_KEY)
	topic := os.Getenv(REPORT_GENERATED_TOPIC_KEY)

	return messagebroker.NewKafkaJsonMessageBroker[JobScheduled[OpenAiJob]](
		logger,
		address,
		topic,
		username,
		password,
	)
}

func ProvideAsJobScheduledBroker(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(messagebroker.MessageBroker[JobScheduled[OpenAiJob]])),
	)
}
