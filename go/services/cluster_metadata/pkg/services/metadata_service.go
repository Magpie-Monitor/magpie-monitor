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
	applicationAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata], nodeAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata],
	clusterAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedClusterState], eventEmitter *EventEmitter) *MetadataService {
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
		nodeAggregatedRepo:          nodeAggregatedRepo,
		clusterStateAggregatedRepo:  clusterAggregatedRepo,
		eventEmitter:                eventEmitter,
		clusterActivityWindowMillis: window,
	}
}

type MetadataService struct {
	log                         *zap.Logger
	clusterRepo                 *sharedrepo.MongoDbCollection[repositories.ClusterState]
	nodeRepo                    *sharedrepo.MongoDbCollection[repositories.NodeState]
	applicationAggregatedRepo   *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
	nodeAggregatedRepo          *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
	clusterStateAggregatedRepo  *sharedrepo.MongoDbCollection[repositories.AggregatedClusterState]
	eventEmitter                *EventEmitter
	clusterActivityWindowMillis int64
}

func (m *MetadataService) InsertClusterMetadata(metadata repositories.ClusterState) error {
	_, err := m.clusterRepo.InsertDocuments([]interface{}{metadata})
	if err != nil {
		m.log.Error("Failed to insert cluster metadata", zap.Error(err))
		return err
	}

	err = m.updateClusterAggregatedState()
	if err != nil {
		m.log.Info(err.Error())
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
	if err != nil {
		m.log.Error("Failed to insert node metadata", zap.Error(err))
		return err
	}

	fmt.Println("updating cluster aggregated state!")
	err = m.updateClusterAggregatedState()
	if err != nil {
		m.log.Info(err.Error())
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

func (m *MetadataService) updateNodeMetadataState(clusterId string, watchedFiles []string) error {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - 300_000}}}},
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
		nodes = append(nodes, repositories.NodeMetadata{Name: n.(string), Files: fileset, Running: true})
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

func (m *MetadataService) updateClusterAggregatedState() error {
	// compute latest state
	// compare with current state
	// if diff -> push event

	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - 300_000}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: time.Now().UnixMilli()}}}},
		}},
	}

	// get all cluster ID's from last hour
	nodeClusterSet, err := m.nodeRepo.GetDistinctDocumentFieldValues("clusterId", filter)
	if err != nil {
		return err
	}

	count, err := m.clusterStateAggregatedRepo.Count(bson.D{})
	if err != nil {
		return err
	}

	if count == 0 {
		newState := m.createAggregatedClusterState(nodeClusterSet)
		m.clusterStateAggregatedRepo.InsertDocument(newState)
		return nil
	}

	latestState, err := m.clusterStateAggregatedRepo.GetDocument(bson.D{}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		return err
	}

	fmt.Println("latestState: ", latestState.Metadata)
	fmt.Println("nodeClusterSet: ", nodeClusterSet)

	if len(nodeClusterSet) != len(latestState.Metadata) {
		newState := m.createAggregatedClusterState(nodeClusterSet)
		m.clusterStateAggregatedRepo.InsertDocument(newState)
		return err
	}

	clusterSet := make(map[string]struct{}, len(latestState.Metadata))
	for _, metadata := range latestState.Metadata {
		clusterSet[metadata.Name] = struct{}{}
	}

	for _, cluster := range nodeClusterSet {
		_, ok := clusterSet[cluster.(string)]
		if !ok {
			newState := m.createAggregatedClusterState(nodeClusterSet)
			_, err := m.clusterStateAggregatedRepo.InsertDocument(newState)
			return err
		}
	}

	return nil
}

func (m *MetadataService) createAggregatedClusterState(clusterSet []interface{}) repositories.AggregatedClusterState {
	fmt.Println(clusterSet)

	state := make([]repositories.ClusterMetadata, 0, len(clusterSet))
	for _, cluster := range clusterSet {
		state = append(state, repositories.ClusterMetadata{Name: cluster.(string), Running: true})
	}

	metadata := repositories.AggregatedClusterState{CollectedAtMs: time.Now().UnixMilli(), Metadata: state}
	m.eventEmitter.EmitClusterMetadataUpdatedEvent(metadata)

	return metadata
}
