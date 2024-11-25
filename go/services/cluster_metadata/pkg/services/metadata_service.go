package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	sharedrepo "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	CLUSTER_AGGREGATED_STATE_POLL_INTERVAL_ENV_NAME     = "CLUSTER_AGGREGATED_STATE_CHANGE_POLL_INTERVAL_SECONDS"
	NODE_AGGREGATED_STATE_POLL_INTERVAL_ENV_NAME        = "NODE_AGGREGATED_STATE_CHANGE_POLL_INTERVAL_SECONDS"
	APPLICATION_AGGREGATED_STATE_POLL_INTERVAL_ENV_NAME = "APPLICATION_AGGREGATED_STATE_CHANGE_POLL_INTERVAL_SECONDS"
	NODE_ACTIVITY_WINDOW_MILLIS_ENV_NAME                = "NODE_ACTIVITY_WINDOW_MILLIS"
	APPLICATION_ACTIVITY_WINDOW_MILLIS_ENV_NAME         = "APPLICATION_ACTIVITY_WINDOW_MILLIS"
	CLUSTER_ACTIVITY_WINDOW_MILLIS_ENV_NAME             = "CLUSTER_ACTIVITY_WINDOW_MILLIS"
	POD_AGENT_APPLICATION_METADATA_TOPIC_ENV_NAME       = "POD_AGENT_APPLICATION_METADATA_TOPIC"
	NODE_AGENT_METADATA_TOPIC_ENV_NAME                  = "NODE_AGENT_NODE_METADATA_TOPIC"
)

func NewApplicationMetadataBroker(logger *zap.Logger, creds *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[repositories.ApplicationState] {
	envs.ValidateEnvs("%s env variable not set", []string{
		POD_AGENT_APPLICATION_METADATA_TOPIC_ENV_NAME,
	})
	return messagebroker.NewKafkaJsonMessageBroker[repositories.ApplicationState](
		logger,
		creds.Address,
		os.Getenv(POD_AGENT_APPLICATION_METADATA_TOPIC_ENV_NAME),
		creds.Username,
		creds.Password,
	)
}

func NewNodeMetadataBroker(logger *zap.Logger, creds *messagebroker.KafkaCredentials) *messagebroker.KafkaJsonMessageBroker[repositories.NodeState] {
	envs.ValidateEnvs("%s env variable not set", []string{
		NODE_AGENT_METADATA_TOPIC_ENV_NAME,
	})
	return messagebroker.NewKafkaJsonMessageBroker[repositories.NodeState](
		logger,
		creds.Address,
		os.Getenv(NODE_AGENT_METADATA_TOPIC_ENV_NAME),
		creds.Username,
		creds.Password,
	)
}

type MetadataServiceParams struct {
	fx.In
	Lc                        fx.Lifecycle
	Logger                    *zap.Logger
	ClusterRepo               *sharedrepo.MongoDbCollection[repositories.ApplicationState]
	NodeRepo                  *sharedrepo.MongoDbCollection[repositories.NodeState]
	ApplicationAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
	NodeAggregatedRepo        *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
	ClusterAggregatedRepo     *sharedrepo.MongoDbCollection[repositories.AggregatedClusterMetadata]
	EventEmitter              *MetadataEventPublisher
	ApplicationMetadataBroker *messagebroker.KafkaJsonMessageBroker[repositories.ApplicationState]
	NodeMetadataBroker        *messagebroker.KafkaJsonMessageBroker[repositories.NodeState]
}

func NewMetadataService(params MetadataServiceParams) *MetadataService {
	metadataService := MetadataService{
		log:                        params.Logger,
		clusterRepo:                params.ClusterRepo,
		nodeRepo:                   params.NodeRepo,
		applicationAggregatedRepo:  params.ApplicationAggregatedRepo,
		nodeAggregatedRepo:         params.NodeAggregatedRepo,
		clusterStateAggregatedRepo: params.ClusterAggregatedRepo,
		eventEmitter:               params.EventEmitter,
		applicationMetadataBroker:  params.ApplicationMetadataBroker,
		nodeMetadataBroker:         params.NodeMetadataBroker,
		clusterAggregatedStateChangePollIntervalSeconds:     envs.ConvertToInt(CLUSTER_AGGREGATED_STATE_POLL_INTERVAL_ENV_NAME),
		nodeAggregatedStateChangePollIntervalSeconds:        envs.ConvertToInt(NODE_AGGREGATED_STATE_POLL_INTERVAL_ENV_NAME),
		applicationAggregatedStateChangePollIntervalSeconds: envs.ConvertToInt(APPLICATION_AGGREGATED_STATE_POLL_INTERVAL_ENV_NAME),
		nodeActivityWindowMillis:                            envs.ConvertToInt64(NODE_ACTIVITY_WINDOW_MILLIS_ENV_NAME),
		applicationActivityWindowMillis:                     envs.ConvertToInt64(APPLICATION_ACTIVITY_WINDOW_MILLIS_ENV_NAME),
		clusterActivityWindowMillis:                         envs.ConvertToInt64(CLUSTER_ACTIVITY_WINDOW_MILLIS_ENV_NAME),
	}

	params.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go metadataService.pollForClusterStateChange()
			go metadataService.pollForApplicationStateChange()
			go metadataService.pollForNodeStateChange()
			go metadataService.consumeApplicationMetadata()
			go metadataService.consumeNodeMetadata()
			return nil
		},
	})

	return &metadataService
}

type MetadataService struct {
	Lc                                                  fx.Lifecycle
	log                                                 *zap.Logger
	clusterRepo                                         *sharedrepo.MongoDbCollection[repositories.ApplicationState]
	nodeRepo                                            *sharedrepo.MongoDbCollection[repositories.NodeState]
	applicationAggregatedRepo                           *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
	nodeAggregatedRepo                                  *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
	clusterStateAggregatedRepo                          *sharedrepo.MongoDbCollection[repositories.AggregatedClusterMetadata]
	eventEmitter                                        *MetadataEventPublisher
	applicationActivityWindowMillis                     int64 // application is considered as running if it has reported in this period
	nodeActivityWindowMillis                            int64 // cluster node is considered as running if it has reported in this period
	clusterActivityWindowMillis                         int64 // cluster is considered as running if it has reported in this period
	clusterAggregatedStateChangePollIntervalSeconds     int
	nodeAggregatedStateChangePollIntervalSeconds        int
	applicationAggregatedStateChangePollIntervalSeconds int
	applicationMetadataBroker                           *messagebroker.KafkaJsonMessageBroker[repositories.ApplicationState]
	nodeMetadataBroker                                  *messagebroker.KafkaJsonMessageBroker[repositories.NodeState]
}

func (m *MetadataService) consumeApplicationMetadata() {
	var (
		msg = make(chan repositories.ApplicationState)
		err = make(chan error)
	)

	defer close(msg)
	defer close(err)

	go m.applicationMetadataBroker.Subscribe(context.TODO(), msg, err)

	for {
		select {
		case metadata := <-msg:
			m.clusterRepo.InsertDocuments([]interface{}{metadata})
		case error := <-err:
			m.log.Error("Error consuming application metadata", zap.Error(error))
		}
	}
}

func (m *MetadataService) consumeNodeMetadata() {
	var (
		msg = make(chan repositories.NodeState)
		err = make(chan error)
	)

	defer close(msg)
	defer close(err)

	go m.nodeMetadataBroker.Subscribe(context.TODO(), msg, err)

	for {
		select {
		case metadata := <-msg:
			m.nodeRepo.InsertDocuments([]interface{}{metadata})
		case error := <-err:
			m.log.Error("Error consuming node metadata", zap.Error(error))
		}
	}
}

func (m *MetadataService) pollForNodeStateChange() {
	for {
		time.Sleep(time.Duration(m.nodeAggregatedStateChangePollIntervalSeconds) * time.Second)

		var wg sync.WaitGroup

		clusterIds, _ := m.nodeRepo.GetDistinctDocumentFieldValues("clusterId", bson.D{})
		for _, id := range clusterIds {
			wg.Add(1)

			go func(clusterId string) {
				defer wg.Done()

				m.log.Info("Updating node metadata", zap.Any("clusterId", clusterId))
				m.updateNodeMetadataStateForCluster(clusterId)
			}(id.(string))
		}

		wg.Wait()
	}
}

func (m *MetadataService) pollForApplicationStateChange() {
	for {
		time.Sleep(time.Duration(m.applicationAggregatedStateChangePollIntervalSeconds) * time.Second)

		var wg sync.WaitGroup

		clusterIds, _ := m.clusterRepo.GetDistinctDocumentFieldValues("clusterId", bson.D{})
		for _, id := range clusterIds {
			wg.Add(1)

			go func(clusterId string) {
				defer wg.Done()

				m.log.Info("Updating application metadata", zap.Any("clusterId", clusterId))
				m.updateApplicationMetadataStateForCluster(clusterId)
			}(id.(string))
		}

		wg.Wait()
	}
}

func (m *MetadataService) pollForClusterStateChange() {
	for {
		time.Sleep(time.Duration(m.clusterAggregatedStateChangePollIntervalSeconds) * time.Second)

		m.log.Info("Updating cluster aggregated state")

		err := m.updateClusterAggregatedState()
		if err != nil {
			m.log.Error("Error updating cluster aggregated state", zap.Error(err))
		}
	}
}

func (m *MetadataService) updateApplicationMetadataStateForCluster(clusterId string) error {
	m.log.Info("Updating application metadata state", zap.String("clusterId", clusterId))

	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - m.applicationActivityWindowMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: time.Now().UnixMilli()}}}},
			bson.D{{Key: "clusterId", Value: bson.D{{Key: "$eq", Value: clusterId}}}},
		}},
	}

	applicationMetadata, err := m.clusterRepo.GetDocuments(filter, bson.D{})
	if err != nil {
		return err
	}
	applicationSet := m.getMetadataApplicationSet(applicationMetadata)

	count, err := m.applicationAggregatedRepo.Count(bson.D{{Key: "clusterId", Value: clusterId}})
	if err != nil {
		m.log.Error("Error fetching document count for application aggregated collection", zap.Error(err))
		return err
	}

	if count == 0 {
		return m.generateAggregatedApplicationStateForCluster(clusterId, applicationSet)
	}

	latestMetadata, err := m.applicationAggregatedRepo.GetDocument(bson.D{{Key: "clusterId", Value: clusterId}}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		return err
	}

	aggregatedApplicationSet := m.getAggregatedMetadataApplicationSet(latestMetadata.Metadata)

	if m.applicationStateHasChanged(applicationSet, aggregatedApplicationSet) {
		return m.generateAggregatedApplicationStateForCluster(clusterId, applicationSet)
	}

	return nil
}

func (m *MetadataService) getMetadataApplicationSet(appMetadata []repositories.ApplicationState) map[string]repositories.ApplicationMetadata {
	applicationSet := make(map[string]repositories.ApplicationMetadata, 0)
	for _, metadata := range appMetadata {
		for _, app := range metadata.Applications {
			key := fmt.Sprintf("%s-%s", app.Name, app.Kind)
			applicationSet[key] = repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind}
		}
	}

	return applicationSet
}

func (m *MetadataService) getAggregatedMetadataApplicationSet(appMetadata []repositories.ApplicationMetadata) map[string]repositories.ApplicationMetadata {
	aggregatedApplicationSet := make(map[string]repositories.ApplicationMetadata, 0)
	for _, app := range appMetadata {
		key := fmt.Sprintf("%s-%s", app.Name, app.Kind)
		aggregatedApplicationSet[key] = repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind}
	}
	return aggregatedApplicationSet
}

func (m *MetadataService) applicationStateHasChanged(applicationSet map[string]repositories.ApplicationMetadata, aggregatedApplicationSet map[string]repositories.ApplicationMetadata) bool {
	if len(applicationSet) != len(aggregatedApplicationSet) {
		return true
	}

	for app, _ := range applicationSet {
		_, exists := aggregatedApplicationSet[app]
		if !exists {
			return false
		}
	}

	return false
}

func (m *MetadataService) generateAggregatedApplicationStateForCluster(clusterId string, applicationSet map[string]repositories.ApplicationMetadata) error {
	m.log.Info("Application metadata aggregated state has changed, performing an update")

	newState := repositories.AggregatedApplicationMetadata{ClusterId: clusterId}
	newState.CollectedAtMs = time.Now().UnixMilli()

	for _, app := range applicationSet {
		newState.Metadata = append(newState.Metadata, repositories.ApplicationMetadata{Name: app.Name, Kind: app.Kind})
	}

	m.log.Info("Updated state", zap.Any("newState", newState))

	_, err := m.applicationAggregatedRepo.InsertDocument(newState)
	if err != nil {
		m.log.Error("Error inserting updated application metadata", zap.Error(err))
		return err
	}

	err = m.eventEmitter.PublishApplicationMetadataUpdatedEvent(newState)
	if err != nil {
		m.log.Error("Error emitting application metadata updated event", zap.Error(err))
		return err
	}

	return nil
}

func (m *MetadataService) updateNodeMetadataStateForCluster(clusterId string) error {
	m.log.Info("Updating metadata state", zap.String("clusterId", clusterId))

	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - m.nodeActivityWindowMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: time.Now().UnixMilli()}}}},
			bson.D{{Key: "clusterId", Value: bson.D{{Key: "$eq", Value: clusterId}}}},
		}},
	}

	nodes, err := m.nodeRepo.GetDocuments(filter, bson.D{})
	if err != nil {
		m.log.Info("Error fetching nodes", zap.Error(err))
		return err
	}

	aggregatedStateCount, err := m.nodeAggregatedRepo.Count(bson.D{{Key: "clusterId", Value: clusterId}})
	if err != nil {
		m.log.Info("Error fetching node aggregated state count", zap.Error(err))
		return err
	}

	if aggregatedStateCount == 0 {
		m.log.Info("There's no aggregated node state saved, generating the first one")
		return m.generateAggregatedNodeStateForCluster(clusterId)
	}

	latestAggregatedState, err := m.nodeAggregatedRepo.GetDocument(bson.D{{Key: "clusterId", Value: clusterId}}, bson.D{{Key: "collectedAtMs", Value: -1}})
	if err != nil {
		return err
	}

	if m.nodeNamesChanged(nodes, latestAggregatedState) {
		return m.generateAggregatedNodeStateForCluster(clusterId)
	}

	if m.watchedFilesChanged(nodes, latestAggregatedState) {
		return m.generateAggregatedNodeStateForCluster(clusterId)
	}

	return nil
}

func (m *MetadataService) nodeNamesChanged(nodes []repositories.NodeState, nodesAggregated repositories.AggregatedNodeMetadata) bool {
	nodeNames := make(map[string]bool, 0)
	for _, node := range nodes {
		nodeNames[node.NodeName] = true
	}

	aggregatedNodeNames := make(map[string]bool, 0)
	for _, metadata := range nodesAggregated.Metadata {
		aggregatedNodeNames[metadata.Name] = true
	}

	if len(nodeNames) != len(aggregatedNodeNames) {
		return true
	}

	for name, _ := range nodeNames {
		exists := aggregatedNodeNames[name]
		if !exists {
			return true
		}
	}

	return false
}

func (m *MetadataService) watchedFilesChanged(nodes []repositories.NodeState, nodesAggregated repositories.AggregatedNodeMetadata) bool {
	nodeFiles := make(map[string]bool, 0)
	for _, node := range nodes {
		nodeFiles[node.NodeName] = true
	}

	aggregatedFiles := make(map[string]bool, 0)
	for _, node := range nodesAggregated.Metadata {
		aggregatedFiles[node.Name] = true
	}

	if len(nodeFiles) != len(aggregatedFiles) {
		return true
	}

	for file, _ := range nodeFiles {
		exists := aggregatedFiles[file]
		if !exists {
			return true
		}
	}

	return false
}

// Node is considered as running if it has reported state in the nodeActivityWindowMillis period
func (m *MetadataService) generateAggregatedNodeStateForCluster(clusterId string) error {
	m.log.Info("Node metadata state change detected")

	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$gte", Value: time.Now().UnixMilli() - m.nodeActivityWindowMillis}}}},
			bson.D{{Key: "collectedAtMs", Value: bson.D{{Key: "$lte", Value: time.Now().UnixMilli()}}}},
			bson.D{{Key: "clusterId", Value: bson.D{{Key: "$eq", Value: clusterId}}}},
		}},
	}

	nodes, err := m.nodeRepo.GetDocuments(filter, bson.D{})
	if err != nil {
		return err
	}

	fileList := m.getDistinctWatchedFilesForNodes(nodes)

	nodeNames, err := m.nodeRepo.GetDistinctDocumentFieldValues("nodeName", filter)
	if err != nil {
		m.log.Error("Error fetching distinct nodeNames", zap.Error(err))
	}

	nodesMetadata := make([]repositories.NodeMetadata, 0, len(nodes))
	for _, name := range nodeNames {
		nodesMetadata = append(nodesMetadata, repositories.NodeMetadata{Name: name.(string), Files: fileList})
	}

	aggregate := repositories.AggregatedNodeMetadata{ClusterId: clusterId, CollectedAtMs: time.Now().UnixMilli(), Metadata: nodesMetadata}
	_, err = m.nodeAggregatedRepo.InsertDocument(aggregate)
	if err != nil {
		m.log.Error("Error inserting updated application metadata", zap.Error(err))
		return err
	}

	err = m.eventEmitter.PublishNodeMetadataUpdatedEvent(aggregate)
	if err != nil {
		m.log.Error("Error emitting node metadata updated event", zap.Error(err))
		return err
	}

	return nil
}

func (m *MetadataService) getDistinctWatchedFilesForNodes(nodes []repositories.NodeState) []interface{} {
	fileSet := make(map[string]struct{}, 0)
	for _, node := range nodes {
		for _, file := range node.WatchedFiles {
			fileSet[file] = struct{}{}
		}
	}

	fileList := make([]interface{}, 0, len(fileSet))
	for file, _ := range fileSet {
		fileList = append(fileList, file)
	}

	return fileList
}

// Fetches clusters that reported state in last clusterActivityWindowMillis
// If fetched state differs from the latest recorded state, a state update event is emitted
func (m *MetadataService) updateClusterAggregatedState() error {
	currentClusterSet, err := m.getUniqueClusterIdsForPeriod(m.clusterActivityWindowMillis)
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
		err := m.generateClusterAggregatedState(currentClusterSet)
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

	latestClusterSet := make(map[string]struct{}, len(latestState.Metadata))
	for _, cluster := range latestState.Metadata {
		latestClusterSet[cluster.ClusterId] = struct{}{}
	}

	if m.clusterStateHasChanged(currentClusterSet, latestClusterSet) {
		err := m.generateClusterAggregatedState(currentClusterSet)
		if err != nil {
			m.log.Info("Error creating cluster aggregated state", zap.Error(err))
			return err
		}
	}

	return nil
}

func (m *MetadataService) clusterStateHasChanged(currentClusterSet map[string]struct{}, latestClusterSet map[string]struct{}) bool {
	if len(currentClusterSet) != len(latestClusterSet) {
		return true
	}

	for cluster, _ := range currentClusterSet {
		_, exists := latestClusterSet[cluster]
		if !exists {
			return true
		}
	}

	return false
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

func (m *MetadataService) generateClusterAggregatedState(clusterSet map[string]struct{}) error {
	state := make([]repositories.ClusterMetadata, 0, len(clusterSet))
	for cluster, _ := range clusterSet {
		state = append(state, repositories.ClusterMetadata{ClusterId: cluster})
	}

	metadata := repositories.AggregatedClusterMetadata{CollectedAtMs: time.Now().UnixMilli(), Metadata: state}

	_, err := m.clusterStateAggregatedRepo.InsertDocument(metadata)
	if err != nil {
		m.log.Error("Error inserting cluster aggregated state", zap.Error(err))
		return err
	}

	err = m.eventEmitter.PublishClusterMetadataUpdatedEvent(metadata)
	if err != nil {
		m.log.Error("Error emitting cluster metadata updated event", zap.Error(err))
		return err
	}

	return nil
}
