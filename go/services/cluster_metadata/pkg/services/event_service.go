package services

import (
	"os"

	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventEmitter(log *zap.Logger, credentials *messagebroker.KafkaCredentials) *EventEmitter {
	return &EventEmitter{
		log:                       log,
		applicationMetadataWriter: NewApplicationMetadataBroker(log, credentials),
		nodeMetadataWriter:        NewNodeMetadataBroker(log, credentials),
		clusterMetadataWriter:     NewClusterMetadataBroker(log, credentials),
	}
}

func NewApplicationMetadataBroker(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[ApplicationMetadataUpdated] {
	appTopic, ok := os.LookupEnv("CLUSTER_METADATA_APPLICATION_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_APPLICATION_TOPIC env variable not provided")
	}

	return messagebroker.NewKafkaJsonMessageBroker[ApplicationMetadataUpdated](logger, credentials.Address, appTopic, credentials.Username, credentials.Password)
}

func NewNodeMetadataBroker(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[NodeMetadataUpdated] {
	nodeTopic, ok := os.LookupEnv("CLUSTER_METADATA_NODE_TOPIC")
	if !ok {
		panic("CLUSTER_METADATA_NODE_TOPIC env variable not provided")
	}

	return messagebroker.NewKafkaJsonMessageBroker[NodeMetadataUpdated](logger, credentials.Address, nodeTopic, credentials.Username, credentials.Password)
}

func NewClusterMetadataBroker(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[ClusterMetadataUpdated] {
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
	CorrelationId string                                 `json:"correlationId"`
	Metadata      repositories.AggregatedClusterMetadata `json:"metadata"`
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

func (e *EventEmitter) EmitClusterMetadataUpdatedEvent(metadata repositories.AggregatedClusterMetadata) error {
	event := ClusterMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.clusterMetadataWriter.Publish(event.CorrelationId, event)
}
