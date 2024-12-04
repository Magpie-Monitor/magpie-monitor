package brokers

import (
	"context"
	"os"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var APPLICATION_INCIDENT_SOURCED_TOPIC_KEY = "APPLICATION_INCIDENT_SOURCED_BROKER_TOPIC"

type ApplicationIncidentSourcedError = string

const (
	ApplicationIncidentSourcedError_Validation ApplicationIncidentSourcedError = "VALIDATION_ERROR"
	ApplicationIncidentSourcedError_Timeout    ApplicationIncidentSourcedError = "TIMEOUT"
	ApplicationIncidentSourcedError_Internal   ApplicationIncidentSourcedError = "INTERNAL_ERROR"
)

type ApplicationIncidentSourced struct {
	CorrelationId             string                                 `json:"correlationId"`
	ApplicationIncidentSource repositories.ApplicationIncidentSource `json:"applicationIncidentSource"`
	TimestampMs               int64                                  `json:"timestampMs"`
}

func NewApplicationIncidentSourced(applicationIncidentSource *repositories.ApplicationIncidentSource) *ApplicationIncidentSourced {
	return &ApplicationIncidentSourced{
		CorrelationId:             applicationIncidentSource.CorrelationId,
		ApplicationIncidentSource: *applicationIncidentSource,
		TimestampMs:               time.Now().UnixMilli(),
	}
}

func NewApplicationIncidentSourcedBroker(lc fx.Lifecycle, logger *zap.Logger) *messagebroker.KafkaJsonMessageBroker[ApplicationIncidentSourced] {

	envs.ValidateEnvs(
		"address/username/password/topic for ApplicationIncidentSourcedBroker is not set",
		[]string{
			MESSAGE_BROKER_ADDRESS_KEY,
			MESSAGE_BROKER_PASSWORD_KEY,
			MESSAGE_BROKER_USERNAME_KEY,
			APPLICATION_INCIDENT_SOURCED_TOPIC_KEY,
		},
	)

	username := os.Getenv(MESSAGE_BROKER_USERNAME_KEY)
	password := os.Getenv(MESSAGE_BROKER_PASSWORD_KEY)
	address := os.Getenv(MESSAGE_BROKER_ADDRESS_KEY)
	topic := os.Getenv(APPLICATION_INCIDENT_SOURCED_TOPIC_KEY)

	broker := messagebroker.NewKafkaJsonMessageBroker[ApplicationIncidentSourced](
		logger,
		address,
		topic,
		username,
		password,
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing connection to ApplicationIncidentSourced broker")
				err := broker.CloseReader()
				if err != nil {
					logger.Error("Error while disconnecting from ApplicationIncidentSourced broker", zap.Error(err))
				}
				return err
			},
		},
	)

	return broker

}

func ProvideAsApplicationIncidentSourcedBroker(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(messagebroker.MessageBroker[ApplicationIncidentSourced])),
	)
}
