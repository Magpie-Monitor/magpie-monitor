package services

import (
	"fmt"
	"slices"
	"time"

	sharedrepo "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func NewMetadataService(log *zap.Logger, clusterRepo *sharedrepo.MongoDbCollection[repositories.ClusterState], nodeRepo *sharedrepo.MongoDbCollection[repositories.NodeState]) *MetadataService {
	return &MetadataService{log: log, clusterRepo: clusterRepo, nodeRepo: nodeRepo}
}

type ApplicationMetadata struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Running bool   `json:"running"`
}

type NodeMetadata struct {
	Name    string        `json:"name"`
	Running bool          `json:"running"`
	Files   []interface{} `json:"files"`
}

type MetadataService struct {
	log         *zap.Logger
	clusterRepo *sharedrepo.MongoDbCollection[repositories.ClusterState]
	nodeRepo    *sharedrepo.MongoDbCollection[repositories.NodeState]
}

func (m *MetadataService) GetClusterMetadataForTimerange(clusterId string, sinceMillis int, toMillis int) ([]ApplicationMetadata, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			bson.D{{Key: "clusterId", Value: bson.D{{Key: "$eq", Value: clusterId}}}},
		}},
	}

	metadata, err := m.clusterRepo.GetDocuments(filter, bson.D{})
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

	apps := make([]ApplicationMetadata, 0, len(applicationSet))
	for _, v := range applicationSet {
		apps = append(apps, v)
	}

	return apps, nil
}

func (m *MetadataService) GetNodeMetadataForTimerange(clusterId string, sinceMillis int, toMillis int) ([]NodeMetadata, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			bson.D{{Key: "clusterId", Value: bson.D{{Key: "$eq", Value: clusterId}}}},
		}},
	}

	fileSet, err := m.nodeRepo.GetDistinctDocumentFieldValues("watchedFiles", filter)
	if err != nil {
		m.log.Error("Error fetching node metadata:", zap.Error(err))
		return nil, err
	}

	nodeSet, err := m.nodeRepo.GetDistinctDocumentFieldValues("nodeName", filter)
	if err != nil {
		m.log.Error("Error fetching node metadata:", zap.Error(err))
		return nil, err
	}

	// a node is considered running if it reported metadata within 10 minutes
	running, err := m.nodeRepo.GetDistinctDocumentFieldValues("nodeName", bson.D{{Key: "collectedAtMs", Value: time.Now().UnixMilli() - 600_000}})
	if err != nil {
		m.log.Error("Error fetching node metadata:", zap.Error(err))
		return nil, err
	}

	nodes := make([]NodeMetadata, 0, len(nodeSet))
	for _, n := range nodeSet {
		running := slices.Contains(running, n)
		nodes = append(nodes, NodeMetadata{Name: n.(string), Files: fileSet, Running: running})
	}

	return nodes, err
}

func (m *MetadataService) InsertClusterMetadata(metadata repositories.ClusterState) error {
	_, err := m.clusterRepo.InsertDocuments([]interface{}{metadata})
	return err
}

func (m *MetadataService) InsertNodeMetadata(metadata repositories.NodeState) error {
	_, err := m.nodeRepo.InsertDocuments([]interface{}{metadata})
	return err
}
