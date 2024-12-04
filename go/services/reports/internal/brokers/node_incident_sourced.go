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

var NODE_INCIDENT_SOURCED_TOPIC_KEY = "NODE_INCIDENT_SOURCED_BROKER_TOPIC"

type NodeIncidentSourcedError = string

const (
	NodeIncidentSourcedError_Validation NodeIncidentSourcedError = "VALIDATION_ERROR"
	NodeIncidentSourcedError_Timeout    NodeIncidentSourcedError = "TIMEOUT"
	NodeIncidentSourcedError_Internal   NodeIncidentSourcedError = "INTERNAL_ERROR"
)

type NodeIncidentSourced struct {
	CorrelationId      string                           `json:"correlationId"`
	NodeIncidentSource *repositories.NodeIncidentSource `json:"nodeIncidentSource"`
	TimestampMs        int64                            `json:"timestampMs"`
}

func NewNodeIncidentSourced(nodeIncidentSource *repositories.NodeIncidentSource) *NodeIncidentSourced {
	return &NodeIncidentSourced{
		CorrelationId:      nodeIncidentSource.CorrelationId,
		NodeIncidentSource: nodeIncidentSource,
		TimestampMs:        time.Now().UnixMilli(),
	}
}

func NewNodeIncidentSourcedBroker(lc fx.Lifecycle, logger *zap.Logger) *messagebroker.KafkaJsonMessageBroker[NodeIncidentSourced] {

	envs.ValidateEnvs(
		"address/username/password/topic for NodeIncidentSourcedBroker is not set",
		[]string{
			MESSAGE_BROKER_ADDRESS_KEY,
			MESSAGE_BROKER_PASSWORD_KEY,
			MESSAGE_BROKER_USERNAME_KEY,
			NODE_INCIDENT_SOURCED_TOPIC_KEY,
		},
	)

	username := os.Getenv(MESSAGE_BROKER_USERNAME_KEY)
	password := os.Getenv(MESSAGE_BROKER_PASSWORD_KEY)
	address := os.Getenv(MESSAGE_BROKER_ADDRESS_KEY)
	topic := os.Getenv(NODE_INCIDENT_SOURCED_TOPIC_KEY)

	broker := messagebroker.NewKafkaJsonMessageBroker[NodeIncidentSourced](
		logger,
		address,
		topic,
		username,
		password,
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing connection to NodeIncidentSourced broker")
				err := broker.CloseReader()
				if err != nil {
					logger.Error("Error while disconnecting from NodeIncidentSourced broker", zap.Error(err))
				}
				return err
			},
		},
	)

	return broker

}

func ProvideAsNodeIncidentSourcedBroker(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(messagebroker.MessageBroker[NodeIncidentSourced])),
	)
}
