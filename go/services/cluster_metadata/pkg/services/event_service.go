package services

import (
	"encoding/json"
	"fmt"
	"os"

	sharedkafka "github.com/Magpie-Monitor/magpie-monitor/pkg/kafka"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventEmitter(log *zap.Logger, credentials *sharedkafka.KafkaCredentials) *EventEmitter {
	return &EventEmitter{log: log, applicationMetadataWriter: NewApplicationMetadataStreamWriter(credentials), nodeMetadataWriter: NewNodeMetadataStreamWriter(credentials)}
}

func NewApplicationMetadataStreamWriter(credentials *sharedkafka.KafkaCredentials) *sharedkafka.StreamWriter {
	appTopic, ok := os.LookupEnv("CLUSTER_METADATA_APPLICATION_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_APPLICATION_TOPIC env variable not provided")
	}

	return sharedkafka.NewStreamWriter(credentials, appTopic, 0)
}

func NewNodeMetadataStreamWriter(credentials *sharedkafka.KafkaCredentials) *sharedkafka.StreamWriter {
	nodeTopic, ok := os.LookupEnv("CLUSTER_METADATA_NODE_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_NODE_TOPIC env variable not provided")
	}

	return sharedkafka.NewStreamWriter(credentials, nodeTopic, 0)
}

type ApplicationMetadataUpdated struct {
	RequestId string                                     `json:"requestId"`
	Metadata  repositories.AggregatedApplicationMetadata `json:"metadata"`
}

type NodeMetadataUpdated struct {
	RequestId string                              `json:"requestId"`
	Metadata  repositories.AggregatedNodeMetadata `json:"metadata"`
}

type EventEmitter struct {
	log                       *zap.Logger
	applicationMetadataWriter *sharedkafka.StreamWriter
	nodeMetadataWriter        *sharedkafka.StreamWriter
}

func (e *EventEmitter) EmitApplicationMetadataUpdatedEvent(metadata repositories.AggregatedApplicationMetadata) error {
	event := ApplicationMetadataUpdated{RequestId: uuid.New().String(), Metadata: metadata}
	return e.emitEvent(event, e.applicationMetadataWriter)
}

func (e *EventEmitter) EmitNodeMetadataUpdatedEvent(metadata repositories.AggregatedNodeMetadata) error {
	fmt.Println("metadata emitted:", metadata)
	event := NodeMetadataUpdated{RequestId: uuid.New().String(), Metadata: metadata}
	return e.emitEvent(event, e.nodeMetadataWriter)
}

func (e *EventEmitter) emitEvent(event interface{}, writer *sharedkafka.StreamWriter) error {
	jsonEvent, err := json.Marshal(&event)
	if err != nil {
		e.log.Error("Error converting metadata event to JSON", zap.Error(err))
		return err
	}

	fmt.Println("jsonEvent:", string(jsonEvent))

	writer.Write(string(jsonEvent))
	return nil
}
