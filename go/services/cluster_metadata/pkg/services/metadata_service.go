package services

import (
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func NewMetadataService(log *zap.Logger, clusterRepo *repositories.MongoDbCollection[repositories.ClusterState], nodeRepo *repositories.MongoDbCollection[repositories.NodeState]) *MetadataService {
	return &MetadataService{log: log, clusterRepo: clusterRepo, nodeRepo: nodeRepo}
}

type MetadataService struct {
	log         *zap.Logger
	clusterRepo *repositories.MongoDbCollection[repositories.ClusterState]
	nodeRepo    *repositories.MongoDbCollection[repositories.NodeState]
}

func (m *MetadataService) GetClusterMetadataForTimerange(clusterName string, sinceMillis int, toMillis int) ([]repositories.ClusterState, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			bson.D{{Key: "clusterName", Value: bson.D{{Key: "$eq", Value: clusterName}}}},
		}},
	}

	return m.clusterRepo.GetFilteredDocuments(filter)
}

func (m *MetadataService) GetNodeMetadataForTimerange(nodeName string, sinceMillis int, toMillis int) ([]repositories.NodeState, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			bson.D{{Key: "nodeName", Value: bson.D{{Key: "$eq", Value: nodeName}}}},
		}},
	}

	return m.nodeRepo.GetFilteredDocuments(filter)
}

func (m *MetadataService) InsertClusterMetadata(metadata repositories.ClusterState) error {
	return m.clusterRepo.InsertDocuments([]interface{}{metadata})
}

func (m *MetadataService) InsertNodeMetadata(metadata repositories.NodeState) error {
	return m.nodeRepo.InsertDocuments([]interface{}{metadata})
}
