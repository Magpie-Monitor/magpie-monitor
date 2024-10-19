package services

import (
	"encoding/json"
	"os"

	kafka "github.com/Magpie-Monitor/magpie-monitor/pkg/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventEmitter(log *zap.Logger, appWriter, nodeWriter *kafka.StreamWriter) *EventEmitter {
	return &EventEmitter{log: log, applicationMetadataWriter: appWriter, nodeMetadataWriter: nodeWriter}
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

type ApplicationMetadataRequested struct {
}

type NodeMetadataRequested struct {
}

type ApplicationMetadataCollected struct {
	requestId string
	metadata  []ApplicationMetadata
}

type NodeMetadataCollected struct {
	requestId string
	metadata  []NodeMetadata
}

type EventEmitter struct {
	log                       *zap.Logger
	applicationMetadataWriter *kafka.StreamWriter
	nodeMetadataWriter        *kafka.StreamWriter
}

func (e *EventEmitter) EmitApplicationMetadataEvent(metadata []ApplicationMetadata) error {
	event := ApplicationMetadataCollected{requestId: uuid.New().String(), metadata: metadata}
	return e.emitEvent(event, e.applicationMetadataWriter)
}

func (e *EventEmitter) EmitNodeMetadataEvent(metadata []NodeMetadata) error {
	event := NodeMetadataCollected{requestId: uuid.New().String(), metadata: metadata}
	return e.emitEvent(event, e.nodeMetadataWriter)
}

func (e *EventEmitter) emitEvent(event interface{}, writer *kafka.StreamWriter) error {
	jsonEvent, err := json.Marshal(&event)
	if err != nil {
		e.log.Error("Error converting metadata event to JSON", zap.Error(err))
		return err
	}

	writer.Write(string(jsonEvent))
	return nil
}
