package services

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	sharedrepo "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func NewMetadataService(log *zap.Logger, clusterRepo *sharedrepo.MongoDbCollection[repositories.ClusterState], nodeRepo *sharedrepo.MongoDbCollection[repositories.NodeState],
	applicationAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata], eventEmitter *EventEmitter) *MetadataService {
	clusterActivityWindowMillis, present := os.LookupEnv("CLUSTER_METADATA_SERVICE_CLUSTER_ACTIVITY_WINDOW_MILLIS")
	if !present {
		panic("env variable CLUSTER_METADATA_SERVICE_CLUSTER_ACTIVITY_WINDOW_MILLIS not set")
	}

	window, err := strconv.ParseInt(clusterActivityWindowMillis, 10, 64)
	if err != nil {
		panic("invalid value for env variable CLUSTER_METADATA_SERVICE_CLUSTER_ACTIVITY_WINDOW_MILLIS, please make sure it's numeric")
	}

	return &MetadataService{
		log:                         log,
		clusterRepo:                 clusterRepo,
		nodeRepo:                    nodeRepo,
		applicationAggregatedRepo:   applicationAggregatedRepo,
		eventEmitter:                eventEmitter,
		clusterActivityWindowMillis: window,
	}
}

type MetadataService struct {
	log                         *zap.Logger
	clusterRepo                 *sharedrepo.MongoDbCollection[repositories.ClusterState]
	nodeRepo                    *sharedrepo.MongoDbCollection[repositories.NodeState]
	applicationAggregatedRepo   *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
	eventEmitter                *EventEmitter
	clusterActivityWindowMillis int64
}

func (m *MetadataService) GetClusterList() ([]repositories.ClusterMetadata, error) {
	clusters, err := m.clusterRepo.GetDistinctDocumentFieldValues("clusterName", bson.D{})
	if err != nil {
		m.log.Error("Error fetching cluster list:", zap.Error(err))
		return nil, err
	}

	// a cluster is considered running if it reported state in the last hour
	toMillis := time.Now().UnixMilli()
	sinceMillis := toMillis - m.clusterActivityWindowMillis

	activeClusters, err := m.clusterRepo.GetDistinctDocumentFieldValues("clusterName",
		bson.D{
			{Key: "$and", Value: bson.A{
				bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: sinceMillis}}}},
				bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: toMillis}}}},
			}},
		})
	if err != nil {
		m.log.Error("Error fetching active cluster list:", zap.Error(err))
		return nil, err
	}

	activeClusterSet := map[string]struct{}{}
	for _, c := range activeClusters {
		clusterName := c.(string)
		activeClusterSet[clusterName] = struct{}{}
	}

	clusterMetadata := make([]repositories.ClusterMetadata, 0)
	for _, c := range clusters {
		clusterName := c.(string)
		_, running := activeClusterSet[clusterName]
		clusterMetadata = append(clusterMetadata, repositories.ClusterMetadata{Name: clusterName, Running: running})
	}

	return clusterMetadata, nil
}

func (m *MetadataService) GetClusterMetadataForTimerange(clusterId string, sinceMillis int, toMillis int) ([]repositories.ApplicationMetadata, error) {
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

	applicationSet := map[string]repositories.ApplicationMetadata{}
	for _, md := range metadata {
		for _, app := range md.Applications {
			// TODO - ADD NAMESPACE
			key := fmt.Sprintf("%s-%s", app.Name, app.Kind)
			_, ok := applicationSet[key]
			if !ok {
				applicationSet[key] = repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind, Running: false}
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

	apps := make([]repositories.ApplicationMetadata, 0, len(applicationSet))
	for _, v := range applicationSet {
		apps = append(apps, v)
	}

	return apps, nil
}

func (m *MetadataService) GetNodeMetadataForTimerange(clusterId string, sinceMillis int, toMillis int) ([]repositories.NodeMetadata, error) {
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

	nodes := make([]repositories.NodeMetadata, 0, len(nodeSet))
	for _, n := range nodeSet {
		running := slices.Contains(running, n)
		nodes = append(nodes, repositories.NodeMetadata{Name: n.(string), Files: fileSet, Running: running})
	}

	return nodes, err
}

func (m *MetadataService) InsertClusterMetadata(metadata repositories.ClusterState) error {
	_, err := m.clusterRepo.InsertDocuments([]interface{}{metadata})
	if err != nil {
		m.log.Error("Failed to insert cluster metadata", zap.Error(err))
		return err
	}

	applicationSet := make(map[string]repositories.Application, 0)
	for _, app := range metadata.Applications {
		applicationSet[app.Name] = app
	}

	count, err := m.applicationAggregatedRepo.Count(bson.D{{Key: "clusterId", Value: metadata.ClusterId}})

	m.log.Info("Aggregated application state count:", zap.Int64("count", count))

	if err != nil {
		m.log.Error("Error fetching document count for application aggregated collection", zap.Error(err))
		return err
	}

	if count == 0 {
		return m.updateApplicationMetadataState(metadata.ClusterId, applicationSet)
	}

	latestState, err := m.applicationAggregatedRepo.GetDocument(bson.D{{Key: "clusterId", Value: metadata.ClusterId}}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		m.log.Error("Error fetching aggregated application metadata", zap.Error(err))
		return err
	}

	m.log.Info("Latest aggregated application state: ", zap.Any("state", latestState))

	if len(applicationSet) != len(latestState.Metadata) {
		return m.updateApplicationMetadataState(metadata.ClusterId, applicationSet)
	}

	for _, app := range latestState.Metadata {
		_, ok := applicationSet[app.Name]
		if !ok {
			return m.updateApplicationMetadataState(metadata.ClusterId, applicationSet)
		}
	}

	for _, app := range applicationSet {
		ok := slices.Contains(latestState.Metadata, repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind, Running: true})
		if !ok {
			return m.updateApplicationMetadataState(metadata.ClusterId, applicationSet)
		}
	}

	return nil
}

func (m *MetadataService) updateApplicationMetadataState(clusterId string, applicationSet map[string]repositories.Application) error {
	m.log.Info("Application metadata aggregated state has changed, performing an update")

	newState := repositories.AggregatedApplicationMetadata{ClusterId: clusterId}
	newState.CollectedAtMs = time.Now().UnixMilli()

	for _, app := range applicationSet {
		newState.Metadata = append(newState.Metadata, repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind, Running: true})
	}

	_, err := m.applicationAggregatedRepo.InsertDocument(newState)
	if err != nil {
		m.log.Error("Error inserting updated application metadata", zap.Error(err))
		return err
	}

	err = m.eventEmitter.EmitApplicationMetadataUpdatedEvent(newState)
	if err != nil {
		m.log.Error("Error emitting application metadata updated event", zap.Error(err))
		return err
	}

	return nil
}

func (m *MetadataService) InsertNodeMetadata(metadata repositories.NodeState) error {
	_, err := m.nodeRepo.InsertDocuments([]interface{}{metadata})
	return err
}
