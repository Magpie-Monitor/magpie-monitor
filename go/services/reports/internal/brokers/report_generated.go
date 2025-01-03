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

var REPORT_GENERATED_TOPIC_KEY = "REPORT_GENERATED_BROKER_TOPIC"

type ReportGenerated struct {
	CorrelationId string               `json:"correlationId"`
	Report        *repositories.Report `json:"report"`
	TimestampMs   int64                `json:"timestampMs"`
}

func NewReportGenerated(report *repositories.Report) ReportGenerated {

	return ReportGenerated{
		CorrelationId: report.CorrelationId,
		Report:        report,
		TimestampMs:   time.Now().UnixMilli(),
	}

}

func NewReportGeneratedBroker(lc fx.Lifecycle, logger *zap.Logger) *messagebroker.KafkaJsonMessageBroker[ReportGenerated] {

	envs.ValidateEnvs(
		"address/username/password/topic for ReportGeneratedBroker is not set",
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

	broker := messagebroker.NewKafkaJsonMessageBroker[ReportGenerated](
		logger,
		address,
		topic,
		username,
		password,
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing connection to ReportGenerated broker")
				err := broker.CloseReader()
				if err != nil {
					logger.Error("Error while disconnecting from ReportGenerated broker", zap.Error(err))

				}
				return err
			},
		},
	)

	return broker
}

func ProvideAsReportGeneratedBroker(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(messagebroker.MessageBroker[ReportGenerated])),
	)
}
