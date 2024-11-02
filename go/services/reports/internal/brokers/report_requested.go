package brokers

import (
	"context"
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var REPORT_REQUESTED_TOPIC_KEY = "REPORT_REQUESTED_BROKER_TOPIC"

type ReportRequest struct {
	ClusterId                *string                                     `json:"clusterId"`
	SinceMs                  *int64                                      `json:"sinceMs"`
	ToMs                     *int64                                      `json:"toMs"`
	ApplicationConfiguration []*insights.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*insights.NodeInsightConfiguration        `json:"nodeConfiguration"`
	MaxLength                *int                                        `json:"maxLength"`
}

type ReportRequested struct {
	CorrelationId string        `json:"correlationId"`
	ReportRequest ReportRequest `json:"reportRequest"`
}

func NewReportRequestedBroker(lc fx.Lifecycle, logger *zap.Logger) *messagebroker.KafkaJsonMessageBroker[ReportRequested] {

	envs.ValidateEnvs(
		"address/username/password/topic for ReportRequestedBroker is not set",
		[]string{
			MESSAGE_BROKER_ADDRESS_KEY,
			MESSAGE_BROKER_PASSWORD_KEY,
			MESSAGE_BROKER_USERNAME_KEY,
			REPORT_REQUESTED_TOPIC_KEY,
		},
	)

	username := os.Getenv(MESSAGE_BROKER_USERNAME_KEY)
	password := os.Getenv(MESSAGE_BROKER_PASSWORD_KEY)
	address := os.Getenv(MESSAGE_BROKER_ADDRESS_KEY)
	topic := os.Getenv(REPORT_REQUESTED_TOPIC_KEY)

	broker := messagebroker.NewKafkaJsonMessageBroker[ReportRequested](
		logger,
		address,
		topic,
		username,
		password,
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing connection to ReportRequested broker")
				err := broker.CloseReader()
				if err != nil {
					logger.Error("Error while disconnecting from ReportRequested broker", zap.Error(err))
				}
				return err
			},
		},
	)

	return broker
}

func ProvideAsReportRequestedBroker(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(messagebroker.MessageBroker[ReportRequested])),
	)
}
