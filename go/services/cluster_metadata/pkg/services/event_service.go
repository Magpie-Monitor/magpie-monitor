package services

import (
	"encoding/json"
	"os"

	sharedkafka "github.com/Magpie-Monitor/magpie-monitor/pkg/kafka"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventEmitter(log *zap.Logger, credentials *sharedkafka.KafkaCredentials) *EventEmitter {
	return &EventEmitter{
		log:                       log,
		applicationMetadataWriter: NewApplicationMetadataStreamWriter(credentials),
		nodeMetadataWriter:        NewNodeMetadataStreamWriter(credentials),
		clusterMetadataWriter:     NewClusterMetadataStreamWriter(credentials),
	}
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

func NewClusterMetadataStreamWriter(credentials *sharedkafka.KafkaCredentials) *sharedkafka.StreamWriter {
	clusterTopic, ok := os.LookupEnv("CLUSTER_METADATA_CLUSTER_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_NODE_TOPIC env variable not provided")
	}

	return sharedkafka.NewStreamWriter(credentials, clusterTopic, 0)
}

type ApplicationMetadataUpdated struct {
	CorrelationId string                                     `json:"correlationId"`
	Metadata      repositories.AggregatedApplicationMetadata `json:"metadata"`
}

type NodeMetadataUpdated struct {
	CorrelationId string                              `json:"correlationId"`
	Metadata      repositories.AggregatedNodeMetadata `json:"metadata"`
}

type ClusterMetadataUpdated struct {
	CorrelationId string                              `json:"correlationId"`
	Metadata      repositories.AggregatedClusterState `json:"metadata"`
}

type EventEmitter struct {
	log                       *zap.Logger
	applicationMetadataWriter *sharedkafka.StreamWriter
	nodeMetadataWriter        *sharedkafka.StreamWriter
	clusterMetadataWriter     *sharedkafka.StreamWriter
}

func (e *EventEmitter) EmitApplicationMetadataUpdatedEvent(metadata repositories.AggregatedApplicationMetadata) error {
	event := ApplicationMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.emitEvent(event, e.applicationMetadataWriter)
}

func (e *EventEmitter) EmitNodeMetadataUpdatedEvent(metadata repositories.AggregatedNodeMetadata) error {
	event := NodeMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.emitEvent(event, e.nodeMetadataWriter)
}

func (e *EventEmitter) EmitClusterMetadataUpdatedEvent(metadata repositories.AggregatedClusterState) error {
	event := ClusterMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.emitEvent(event, e.clusterMetadataWriter)
}

func (e *EventEmitter) emitEvent(event interface{}, writer *sharedkafka.StreamWriter) error {
	jsonEvent, err := json.Marshal(&event)
	if err != nil {
		e.log.Error("Error converting metadata event to JSON", zap.Error(err))
		return err
	}

	writer.Write(string(jsonEvent))
	return nil
}
