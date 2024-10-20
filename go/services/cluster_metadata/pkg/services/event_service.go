package services

import (
	"encoding/json"
	"os"

	kafka "github.com/Magpie-Monitor/magpie-monitor/pkg/kafka"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventEmitter(log *zap.Logger, credentials *kafka.KafkaCredentials) *EventEmitter {
	return &EventEmitter{log: log, applicationMetadataWriter: NewApplicationMetadataStreamWriter(credentials), nodeMetadataWriter: NewNodeMetadataStreamWriter(credentials)}
}

func NewApplicationMetadataStreamWriter(credentials *kafka.KafkaCredentials) *kafka.StreamWriter {
	appTopic, ok := os.LookupEnv("CLUSTER_METADATA_APPLICATION_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_APPLICATION_TOPIC env variable not provided")
	}

	return kafka.NewStreamWriter(credentials, appTopic, 0)
}

func NewNodeMetadataStreamWriter(credentials *kafka.KafkaCredentials) *kafka.StreamWriter {
	nodeTopic, ok := os.LookupEnv("CLUSTER_METADATA_NODE_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_NODE_TOPIC env variable not provided")
	}

	return kafka.NewStreamWriter(credentials, nodeTopic, 0)
}

type ApplicationMetadataUpdated struct {
	requestId string
	metadata  repositories.AggregatedApplicationMetadata
}

type NodeMetadataUpdated struct {
	requestId string
	metadata  repositories.AggregatedNodeMetadata
}

type EventEmitter struct {
	log                       *zap.Logger
	applicationMetadataWriter *kafka.StreamWriter
	nodeMetadataWriter        *kafka.StreamWriter
}

func (e *EventEmitter) EmitApplicationMetadataUpdatedEvent(metadata repositories.AggregatedApplicationMetadata) error {
	event := ApplicationMetadataUpdated{requestId: uuid.New().String(), metadata: metadata}
	return e.emitEvent(event, e.applicationMetadataWriter)
}

func (e *EventEmitter) EmitNodeMetadataUpdatedEvent(metadata repositories.AggregatedNodeMetadata) error {
	event := NodeMetadataUpdated{requestId: uuid.New().String(), metadata: metadata}
	return e.emitEvent(event, e.nodeMetadataWriter)
}

func (e *EventEmitter) emitEvent(event interface{}, writer *kafka.StreamWriter) error {
	e.log.Info("Emitting an event:", zap.Any("event", event))

	jsonEvent, err := json.Marshal(&event)
	if err != nil {
		e.log.Error("Error converting metadata event to JSON", zap.Error(err))
		return err
	}

	writer.Write(string(jsonEvent))
	return nil
}
