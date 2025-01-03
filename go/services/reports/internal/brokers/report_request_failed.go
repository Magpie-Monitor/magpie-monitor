package brokers

import (
	"context"
	"os"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var REPORT_REQUEST_FAILED_TOPIC_KEY = "REPORT_REQUEST_FAILED_BROKER_TOPIC"

type ReportRequestError = string

const (
	ReportRequestError_Validation ReportRequestError = "VALIDATION_ERROR"
	ReportRequestError_Timeout    ReportRequestError = "TIMEOUT"
	ReportRequestError_Internal   ReportRequestError = "INTERNAL_ERROR"
)

type ReportRequestFailed struct {
	CorrelationId string             `json:"correlationId"`
	ErrorType     ReportRequestError `json:"errorType"`
	ErrorMessage  string             `json:"errorMessage"`
	TimestampMs   int64              `json:"timestampMs"`
}

func NewReportRequestFailed(correlationId string, errorType ReportRequestError, msg string) *ReportRequestFailed {
	return &ReportRequestFailed{
		CorrelationId: correlationId,
		ErrorType:     errorType,
		ErrorMessage:  msg,
		TimestampMs:   time.Now().UnixMilli(),
	}
}

func NewReportRequestFailedValidation(correlationId string, msg string) *ReportRequestFailed {
	return NewReportRequestFailed(
		correlationId,
		ReportRequestError_Validation,
		msg,
	)
}

func NewReportRequestFailedInternalError(correlationId string, msg string) *ReportRequestFailed {
	return NewReportRequestFailed(
		correlationId,
		ReportRequestError_Internal,
		msg,
	)
}

func NewReportRequestFailedBroker(lc fx.Lifecycle, logger *zap.Logger) *messagebroker.KafkaJsonMessageBroker[ReportRequestFailed] {

	envs.ValidateEnvs(
		"address/username/password/topic for ReportRequestFailedBroker is not set",
		[]string{
			MESSAGE_BROKER_ADDRESS_KEY,
			MESSAGE_BROKER_PASSWORD_KEY,
			MESSAGE_BROKER_USERNAME_KEY,
			REPORT_REQUEST_FAILED_TOPIC_KEY,
		},
	)

	username := os.Getenv(MESSAGE_BROKER_USERNAME_KEY)
	password := os.Getenv(MESSAGE_BROKER_PASSWORD_KEY)
	address := os.Getenv(MESSAGE_BROKER_ADDRESS_KEY)
	topic := os.Getenv(REPORT_REQUEST_FAILED_TOPIC_KEY)

	broker := messagebroker.NewKafkaJsonMessageBroker[ReportRequestFailed](
		logger,
		address,
		topic,
		username,
		password,
	)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing connection to ReportRequestFailed broker")
				err := broker.CloseReader()
				if err != nil {
					logger.Error("Error while disconnecting from ReportRequestFailed broker", zap.Error(err))
				}
				return err
			},
		},
	)

	return broker

}

func ProvideAsReportRequestFailedBroker(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(messagebroker.MessageBroker[ReportRequestFailed])),
	)
}
