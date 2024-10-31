package services

import (
	"context"
	"encoding/json"
	"os"

	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventEmitter(log *zap.Logger, credentials *messagebroker.KafkaCredentials) *EventEmitter {
	return &EventEmitter{
		log:                       log,
		applicationMetadataWriter: NewApplicationMetadataStreamWriter(log, credentials),
		nodeMetadataWriter:        NewNodeMetadataStreamWriter(log, credentials),
		clusterMetadataWriter:     NewClusterMetadataStreamWriter(log, credentials),
	}
}

func NewApplicationMetadataStreamWriter(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[ApplicationMetadataUpdated] {
	appTopic, ok := os.LookupEnv("CLUSTER_METADATA_APPLICATION_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_APPLICATION_TOPIC env variable not provided")
	}

	return messagebroker.NewKafkaJsonMessageBroker[ApplicationMetadataUpdated](logger, credentials.Address, appTopic, credentials.Username, credentials.Password)
}

func NewNodeMetadataStreamWriter(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[NodeMetadataUpdated] {
	nodeTopic, ok := os.LookupEnv("CLUSTER_METADATA_NODE_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_NODE_TOPIC env variable not provided")
	}

	return messagebroker.NewKafkaJsonMessageBroker[NodeMetadataUpdated](logger, credentials.Address, nodeTopic, credentials.Username, credentials.Password)
}

func NewClusterMetadataStreamWriter(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[ClusterMetadataUpdated] {
	clusterTopic, ok := os.LookupEnv("CLUSTER_METADATA_CLUSTER_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_NODE_TOPIC env variable not provided")
	}

	return messagebroker.NewKafkaJsonMessageBroker[ClusterMetadataUpdated](logger, credentials.Address, clusterTopic, credentials.Username, credentials.Password)
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
	applicationMetadataWriter *messagebroker.KafkaJsonMessageBroker[ApplicationMetadataUpdated]
	nodeMetadataWriter        *messagebroker.KafkaJsonMessageBroker[NodeMetadataUpdated]
	clusterMetadataWriter     *messagebroker.KafkaJsonMessageBroker[ClusterMetadataUpdated]
}

func (e *EventEmitter) EmitApplicationMetadataUpdatedEvent(metadata repositories.AggregatedApplicationMetadata) error {
	event := ApplicationMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.applicationMetadataWriter.Publish(event.CorrelationId, event)
}

func (e *EventEmitter) EmitNodeMetadataUpdatedEvent(metadata repositories.AggregatedNodeMetadata) error {
	event := NodeMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.nodeMetadataWriter.Publish(event.CorrelationId, event)
}

func (e *EventEmitter) EmitClusterMetadataUpdatedEvent(metadata repositories.AggregatedClusterState) error {
	event := ClusterMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.clusterMetadataWriter.Publish(event.CorrelationId, event)
}

func (e *EventEmitter) emitEvent(event interface{}, correlationId string, writer *messagebroker.KafkaMessageBroker) error {
	jsonEvent, err := json.Marshal(&event)
	if err != nil {
		e.log.Error("Error converting metadata event to JSON", zap.Error(err))
		return err
	}

	writer.Publish(context.Background(), []byte(correlationId), []byte(jsonEvent))

	return nil
}
