package services

import (
	"fmt"
	"log"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func NewMetadataService(log *zap.Logger, clusterRepo *repositories.MongoDbCollection[repositories.ClusterState], nodeRepo *repositories.MongoDbCollection[repositories.NodeState]) *MetadataService {
	return &MetadataService{log: log, clusterRepo: clusterRepo, nodeRepo: nodeRepo}
}

type ApplicationMetadata struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Running bool   `json:"running"`
}

type NodeMetadata struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
	Files   string `json:"files"`
}

type MetadataService struct {
	log         *zap.Logger
	clusterRepo *repositories.MongoDbCollection[repositories.ClusterState]
	nodeRepo    *repositories.MongoDbCollection[repositories.NodeState]
}

func (m *MetadataService) GetClusterMetadataForTimerange(clusterName string, sinceMillis int, toMillis int) ([]ApplicationMetadata, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			bson.D{{Key: "clusterName", Value: bson.D{{Key: "$eq", Value: clusterName}}}},
		}},
	}

	metadata, err := m.clusterRepo.GetDocuments(filter, bson.D{})
	log.Println("metadata:", metadata)
	if err != nil {
		m.log.Error("Error fetching cluster metadata:", zap.Error(err))
		return nil, err
	}

	applicationSet := map[string]ApplicationMetadata{}
	for _, md := range metadata {
		for _, app := range md.Applications {
			// TODO - ADD NAMESPACE
			key := fmt.Sprintf("%s-%s", app.Name, app.Kind)
			_, ok := applicationSet[key]
			if !ok {
				applicationSet[key] = ApplicationMetadata{Name: app.Name, Kind: app.Kind, Running: false}
			}
		}
	}

	running, err := m.clusterRepo.GetDocument(bson.D{}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		m.log.Error("Error fetching cluster metadata:", zap.Error(err))
		return nil, err
	}

	for _, app := range running.Applications {
		key := fmt.Sprintf("%s-%s", app.Name, app.Kind)
		val, ok := applicationSet[key]
		if ok {
			val.Running = true
			applicationSet[key] = val
		}
	}

	apps := make([]ApplicationMetadata, 0)
	for _, v := range applicationSet {
		apps = append(apps, v)
	}

	return apps, nil
}

func (m *MetadataService) GetNodeMetadataForTimerange(nodeName string, sinceMillis int, toMillis int) ([]repositories.NodeState, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			bson.D{{Key: "nodeName", Value: bson.D{{Key: "$eq", Value: nodeName}}}},
		}},
	}

	return m.nodeRepo.GetDocuments(filter, nil)
}

func (m *MetadataService) InsertClusterMetadata(metadata repositories.ClusterState) error {
	return m.clusterRepo.InsertDocuments([]interface{}{metadata})
}

func (m *MetadataService) InsertNodeMetadata(metadata repositories.NodeState) error {
	return m.nodeRepo.InsertDocuments([]interface{}{metadata})
}
