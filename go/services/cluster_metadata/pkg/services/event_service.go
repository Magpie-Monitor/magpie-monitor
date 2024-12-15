package services

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
)

const (
	CLUSTER_METADATA_APPLICATION_TOPIC_ENV_NAME = "CLUSTER_METADATA_APPLICATION_TOPIC"
	CLUSTER_METADATA_NODE_TOPIC_ENV_NAME        = "CLUSTER_METADATA_NODE_TOPIC"
	CLUSTER_METADATA_CLUSTER_TOPIC_ENV_NAME     = "CLUSTER_METADATA_CLUSTER_TOPIC"
)

func NewMetadataEventPublisher(
	log *zap.Logger,
	appMetadataBroker messagebroker.MessageBroker[ApplicationMetadataUpdated],
	nodeMetadataBroker messagebroker.MessageBroker[NodeMetadataUpdated],
	clusterMetadataBroker messagebroker.MessageBroker[ClusterMetadataUpdated],
) *MetadataEventPublisher {
	return &MetadataEventPublisher{
		log:                       log,
		applicationMetadataBroker: appMetadataBroker,
		nodeMetadataBroker:        nodeMetadataBroker,
		clusterMetadataBroker:     clusterMetadataBroker,
	}
}

func NewApplicationMetadataUpdatedBroker(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) messagebroker.MessageBroker[ApplicationMetadataUpdated] {
	envs.ValidateEnvs("%s env variable not set", []string{
		CLUSTER_METADATA_APPLICATION_TOPIC_ENV_NAME,
	})
	return messagebroker.NewKafkaJsonMessageBroker[ApplicationMetadataUpdated](
		logger,
		credentials.Address, os.Getenv(CLUSTER_METADATA_APPLICATION_TOPIC_ENV_NAME),
		credentials.Username,
		credentials.Password,
	)
}

func NewNodeMetadataUpdatedBroker(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) messagebroker.MessageBroker[NodeMetadataUpdated] {
	envs.ValidateEnvs("%s env variable not set", []string{
		CLUSTER_METADATA_NODE_TOPIC_ENV_NAME,
	})
	return messagebroker.NewKafkaJsonMessageBroker[NodeMetadataUpdated](
		logger,
		credentials.Address,
		os.Getenv(CLUSTER_METADATA_NODE_TOPIC_ENV_NAME),
		credentials.Username,
		credentials.Password,
	)
}

func NewClusterMetadataUpdatedBroker(logger *zap.Logger, credentials *messagebroker.KafkaCredentials) messagebroker.MessageBroker[ClusterMetadataUpdated] {
	envs.ValidateEnvs("%s env variable not set", []string{
		CLUSTER_METADATA_CLUSTER_TOPIC_ENV_NAME,
	})
	return messagebroker.NewKafkaJsonMessageBroker[ClusterMetadataUpdated](
		logger,
		credentials.Address,
		os.Getenv(CLUSTER_METADATA_CLUSTER_TOPIC_ENV_NAME),
		credentials.Username,
		credentials.Password,
	)
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

type MetadataEventPublisher struct {
	log                       *zap.Logger
	applicationMetadataBroker messagebroker.MessageBroker[ApplicationMetadataUpdated]
	nodeMetadataBroker        messagebroker.MessageBroker[NodeMetadataUpdated]
	clusterMetadataBroker     messagebroker.MessageBroker[ClusterMetadataUpdated]
}

func (e *MetadataEventPublisher) PublishApplicationMetadataUpdatedEvent(metadata repositories.AggregatedApplicationMetadata) error {
	e.log.Info("Publishing application metadata updated event", zap.Any("event", metadata))
	event := ApplicationMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.applicationMetadataBroker.Publish(event.CorrelationId, event)
}

func (e *MetadataEventPublisher) PublishNodeMetadataUpdatedEvent(metadata repositories.AggregatedNodeMetadata) error {
	e.log.Info("Publishing node metadata updated event", zap.Any("event", metadata))
	event := NodeMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.nodeMetadataBroker.Publish(event.CorrelationId, event)
}

func (e *MetadataEventPublisher) PublishClusterMetadataUpdatedEvent(metadata repositories.AggregatedClusterMetadata) error {
	e.log.Info("Publishing cluster metadata updated event", zap.Any("event", metadata))
	event := ClusterMetadataUpdated{CorrelationId: uuid.New().String(), Metadata: metadata}
	return e.clusterMetadataBroker.Publish(event.CorrelationId, event)
}
