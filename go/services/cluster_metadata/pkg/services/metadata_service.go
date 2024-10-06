package services

import (
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/entity"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func NewMetadataService(log *zap.Logger, clusterRepo *repositories.MongoDbCollection[entity.ClusterState], nodeRepo *repositories.MongoDbCollection[entity.NodeState]) *MetadataService {
	return &MetadataService{log: log, clusterRepo: clusterRepo, nodeRepo: nodeRepo}
}

type MetadataService struct {
	log         *zap.Logger
	clusterRepo *repositories.MongoDbCollection[entity.ClusterState]
	nodeRepo    *repositories.MongoDbCollection[entity.NodeState]
}

func (m *MetadataService) GetClusterMetadataForTimerange(clusterName string, sinceMillis int, toMillis int) ([]entity.ClusterState, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedatms", Value: bson.D{{Key: "$gt", Value: sinceMillis}}}},
			bson.D{{Key: "collectedatms", Value: bson.D{{Key: "$lt", Value: toMillis}}}},
			bson.D{{Key: "clustername", Value: bson.D{{Key: "$eq", Value: clusterName}}}},
		}},
	}

	return m.clusterRepo.GetFilteredDocuments(filter)
}

func (m *MetadataService) GetNodeMetadataForTimerange(nodeName string, sinceMillis int, toMillis int) ([]entity.NodeState, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedatms", Value: bson.D{{Key: "$gt", Value: sinceMillis}}}},
			bson.D{{Key: "collectedatms", Value: bson.D{{Key: "$lt", Value: toMillis}}}},
			bson.D{{Key: "nodename", Value: bson.D{{Key: "$eq", Value: nodeName}}}},
		}},
	}

	return m.nodeRepo.GetFilteredDocuments(filter)
}

func (m *MetadataService) InsertClusterMetadata(metadata entity.ClusterState) error {
	return m.clusterRepo.InsertDocuments([]interface{}{metadata})
}

func (m *MetadataService) InsertNodeMetadata(metadata entity.NodeState) error {
	return m.nodeRepo.InsertDocuments([]interface{}{metadata})
}
