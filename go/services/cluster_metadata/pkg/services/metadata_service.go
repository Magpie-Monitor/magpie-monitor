package services

import (
	"context"
	"slices"
	"time"

	sharedrepo "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	clusterAggregatedStateUpdateSleepSeconds = 30        // how often cluster state is refreshed
	nodeActivityWindowMillis                 = 3_600_000 // cluster node is considered as running if it has reported in this period
	clusterActivityWindowMillis              = 3_600_000 // cluster is considered as running if it has reported in this period
)

func NewMetadataService(lc fx.Lifecycle, log *zap.Logger, clusterRepo *sharedrepo.MongoDbCollection[repositories.ApplicationState], nodeRepo *sharedrepo.MongoDbCollection[repositories.NodeState],
	applicationAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata], nodeAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata],
	clusterAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedClusterMetadata], eventEmitter *EventEmitter) *MetadataService {

	metadataService := MetadataService{
		log:                        log,
		clusterRepo:                clusterRepo,
		nodeRepo:                   nodeRepo,
		applicationAggregatedRepo:  applicationAggregatedRepo,
		nodeAggregatedRepo:         nodeAggregatedRepo,
		clusterStateAggregatedRepo: clusterAggregatedRepo,
		eventEmitter:               eventEmitter,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go metadataService.scheduleClusterStateUpdate()
			return nil
		},
	})

	return &metadataService
}

type MetadataService struct {
	log                        *zap.Logger
	Lc                         fx.Lifecycle
	clusterRepo                *sharedrepo.MongoDbCollection[repositories.ApplicationState]
	nodeRepo                   *sharedrepo.MongoDbCollection[repositories.NodeState]
	applicationAggregatedRepo  *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
	nodeAggregatedRepo         *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
	clusterStateAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedClusterMetadata]
	eventEmitter               *EventEmitter
}

func (m *MetadataService) InsertApplicationMetadata(metadata repositories.ApplicationState) error {
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
		ok := slices.Contains(latestState.Metadata, repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind})
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
		newState.Metadata = append(newState.Metadata, repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind})
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
	if err != nil {
		m.log.Error("Failed to insert node metadata", zap.Error(err))
		return err
	}

	count, err := m.nodeAggregatedRepo.Count(bson.D{{Key: "clusterId", Value: metadata.ClusterId}})
	m.log.Info("Aggregated application state count:", zap.Int64("count", count))
	if err != nil {
		m.log.Error("Error fetching document count for node aggregated collection", zap.Error(err))
		return err
	}

	if count == 0 {
		return m.updateNodeMetadataState(metadata.ClusterId, metadata.WatchedFiles)
	}

	latestState, err := m.nodeAggregatedRepo.GetDocument(bson.D{{Key: "clusterId", Value: metadata.ClusterId}}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		m.log.Error("Error fetching aggregated node metadata", zap.Error(err))
		return err
	}

	nodeSet := make(map[string]repositories.NodeMetadata, 0)
	for _, node := range latestState.Metadata {
		nodeSet[node.Name] = node
	}

	node, exists := nodeSet[metadata.NodeName]
	if !exists || len(node.Files) != len(metadata.WatchedFiles) {
		return m.updateNodeMetadataState(metadata.ClusterId, metadata.WatchedFiles)
	}

	for _, file := range node.Files {
		if !slices.Contains(metadata.WatchedFiles, file.(string)) {
			return m.updateNodeMetadataState(metadata.ClusterId, metadata.WatchedFiles)
		}
	}

	return nil
}

// Node is considered as running if it has reported state in the nodeActivityWindowMillis period
func (m *MetadataService) updateNodeMetadataState(clusterId string, watchedFiles []string) error {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - nodeActivityWindowMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: time.Now().UnixMilli()}}}},
			bson.D{{Key: "clusterId", Value: bson.D{{Key: "$eq", Value: clusterId}}}},
		}},
	}

	fileset := []interface{}{}
	for _, file := range watchedFiles {
		fileset = append(fileset, file)
	}

	nodeSet, err := m.nodeRepo.GetDistinctDocumentFieldValues("nodeName", filter)
	if err != nil {
		m.log.Error("Error fetching node metadata:", zap.Error(err))
		return err
	}

	nodes := make([]repositories.NodeMetadata, 0)
	for _, n := range nodeSet {
		nodes = append(nodes, repositories.NodeMetadata{Name: n.(string), Files: fileset})
	}

	aggregate := repositories.AggregatedNodeMetadata{ClusterId: clusterId, CollectedAtMs: time.Now().UnixMilli(), Metadata: nodes}
	_, err = m.nodeAggregatedRepo.InsertDocument(aggregate)
	if err != nil {
		m.log.Error("Error inserting updated application metadata", zap.Error(err))
		return err
	}

	err = m.eventEmitter.EmitNodeMetadataUpdatedEvent(aggregate)
	if err != nil {
		m.log.Error("Error emitting node metadata updated event", zap.Error(err))
		return err
	}

	return nil
}

func (m *MetadataService) scheduleClusterStateUpdate() {
	for {
		m.log.Info("Updating cluster aggregated state")

		err := m.updateClusterAggregatedState()
		if err != nil {
			m.log.Error("Error updating cluster aggregated state", zap.Error(err))
		}

		time.Sleep(clusterAggregatedStateUpdateSleepSeconds * time.Second)
	}
}

// Fetches clusters that reported state in last clusterActivityWindowMillis
// If fetched state differs from the latest recorded state, a state update event is emitted
func (m *MetadataService) updateClusterAggregatedState() error {
	clusterSet, err := m.getUniqueClusterIdsForPeriod(clusterActivityWindowMillis)
	if err != nil {
		m.log.Error("Error fetching unique cluster ID's", zap.Error(err))
		return nil
	}

	count, err := m.clusterStateAggregatedRepo.Count(bson.D{})
	if err != nil {
		m.log.Info("Error fetching cluster aggregated state", zap.Error(err))
		return err
	}

	if count == 0 {
		_, err := m.createAggregatedClusterState(clusterSet)
		if err != nil {
			m.log.Info("Error creating cluster aggregated state for count=0", zap.Error(err))
			return err
		}

		return nil
	}

	latestState, err := m.clusterStateAggregatedRepo.GetDocument(bson.D{}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		m.log.Info("Error fetching latest cluster aggregated state", zap.Error(err))
		return err
	}

	if len(clusterSet) != len(latestState.Metadata) {
		_, err := m.createAggregatedClusterState(clusterSet)
		if err != nil {
			m.log.Info("Error creating cluster aggregated state", zap.Error(err))
			return err
		}

		return nil
	}

	latestClusterSet := make(map[string]struct{}, len(latestState.Metadata))
	for _, metadata := range latestState.Metadata {
		latestClusterSet[metadata.ClusterId] = struct{}{}
	}

	for cluster, _ := range clusterSet {
		_, ok := clusterSet[cluster]
		if !ok {
			_, err := m.createAggregatedClusterState(clusterSet)
			if err != nil {
				m.log.Info("Error creating cluster aggregated state", zap.Error(err))
				return err
			}

			return nil
		}
	}

	return nil
}

func (m *MetadataService) getUniqueClusterIdsForPeriod(periodMillis int64) (map[string]struct{}, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - periodMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: time.Now().UnixMilli()}}}},
		}},
	}

	nodeClusterSet, err := m.nodeRepo.GetDistinctDocumentFieldValues("clusterId", filter)
	if err != nil {
		m.log.Info("Error fetching clusterId set from node repository", zap.Error(err))
		return nil, err
	}

	applicationClusterSet, err := m.clusterRepo.GetDistinctDocumentFieldValues("clusterId", filter)
	if err != nil {
		m.log.Info("Error fetching clusterId set from application repository", zap.Error(err))
		return nil, err
	}

	clusters := append(nodeClusterSet, applicationClusterSet...)
	clusterSet := make(map[string]struct{}, 0)
	for _, cluster := range clusters {
		clusterSet[cluster.(string)] = struct{}{}
	}

	return clusterSet, nil
}

func (m *MetadataService) createAggregatedClusterState(clusterSet map[string]struct{}) (repositories.AggregatedClusterMetadata, error) {
	state := make([]repositories.ClusterMetadata, 0, len(clusterSet))
	for cluster, _ := range clusterSet {
		state = append(state, repositories.ClusterMetadata{ClusterId: cluster})
	}

	metadata := repositories.AggregatedClusterMetadata{CollectedAtMs: time.Now().UnixMilli(), Metadata: state}

	_, err := m.clusterStateAggregatedRepo.InsertDocument(metadata)
	if err != nil {
		m.log.Error("Error inserting cluster aggregated state", zap.Error(err))
		return repositories.AggregatedClusterMetadata{}, err
	}

	err = m.eventEmitter.EmitClusterMetadataUpdatedEvent(metadata)
	if err != nil {
		m.log.Error("Error emitting cluster metadata updated event", zap.Error(err))
		return repositories.AggregatedClusterMetadata{}, err
	}

	return metadata, nil
}
