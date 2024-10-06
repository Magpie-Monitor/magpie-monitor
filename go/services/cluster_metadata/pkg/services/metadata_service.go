package services

import (
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/internal/entity"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type MetadataService struct {
	log         *zap.Logger
	clusterRepo repositories.MongoDbRepository[entity.ClusterState]
	nodeRepo    repositories.MongoDbRepository[entity.NodeState]
}

func NewMetadataService(log *zap.Logger, clusterRepo repositories.MongoDbRepository[entity.ClusterState], nodeRepo repositories.MongoDbRepository[entity.NodeState]) MetadataService {
	return MetadataService{log: log, clusterRepo: clusterRepo, nodeRepo: nodeRepo}
}

func (m *MetadataService) GetClusterMetadataForTimerange(clusterName string, sinceMillis int64, toMillis int64) ([]entity.ClusterState, error) {
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"collectedatms", bson.D{{"$gt", sinceMillis}}}},
			bson.D{{"collectedatms", bson.D{{"$lt", toMillis}}}},
			bson.D{{"clustername", bson.D{{"$eq", clusterName}}}},
		}},
	}

	return m.clusterRepo.GetFilteredDocuments(filter)
}

func (m *MetadataService) GetNodeMetadataForTimerange(nodeName string, sinceMillis int64, toMillis int64) ([]entity.NodeState, error) {
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"collectedatms", bson.D{{"$gt", sinceMillis}}}},
			bson.D{{"collectedatms", bson.D{{"$lt", toMillis}}}},
			bson.D{{"nodename", bson.D{{"$eq", nodeName}}}},
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
