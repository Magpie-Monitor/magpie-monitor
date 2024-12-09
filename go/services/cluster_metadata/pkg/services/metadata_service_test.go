package services_test

import (
	"context"
	"testing"
	"time"

	messagebroker "github.com/Magpie-Monitor/magpie-monitor/pkg/message-broker"
	sharedrepo "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestNodeMetadataIngestion(t *testing.T) {

	testCases := []struct {
		name     string
		metadata repositories.NodeState
	}{
		{
			name: "Ingest metadata with multiple watched files",
			metadata: repositories.NodeState{
				ClusterId:     "cluster",
				NodeName:      "node",
				CollectedAtMs: time.Now().UnixMilli(),
				WatchedFiles: []string{
					"file1",
					"file2",
					"file3",
					"file4",
				},
			},
		},
		{
			name: "Ingest metadata without watched files",
			metadata: repositories.NodeState{
				ClusterId:     "cluster",
				NodeName:      "node",
				CollectedAtMs: time.Now().UnixMilli(),
				WatchedFiles:  []string{},
			},
		},
	}

	type TestDependencies struct {
		fx.In
		Logger *zap.Logger

		MetadataService *services.MetadataService

		NodeMetadataBroker messagebroker.MessageBroker[repositories.NodeState]

		NodeMetadataRepository *sharedrepo.MongoDbCollection[repositories.NodeState]
		NodeAggregatedRepo     *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
	}

	test := func(dependencies TestDependencies) {
		log := dependencies.Logger

		for _, test := range testCases {
			dependencies.MetadataService.Init()

			var metadata = test.metadata

			log.Info("Executing", zap.String("test", test.name))

			dependencies.NodeMetadataRepository.DeleteAll()
			dependencies.NodeAggregatedRepo.DeleteAll()

			dependencies.NodeMetadataBroker.Publish("", metadata)

			time.Sleep(15 * time.Second)

			result, err := dependencies.NodeMetadataRepository.GetDocument(
				bson.D{
					{Key: "clusterId", Value: metadata.ClusterId},
					{Key: "collectedAtMs", Value: metadata.CollectedAtMs}},
				bson.D{},
			)

			if err != nil {
				log.Error("Error reading node metadata", zap.Error(err))
				t.Error("Error reading node metadata")
			}

			assert.Len(t, result.WatchedFiles, len(metadata.WatchedFiles), "Invalid number of watched files")
			assert.ElementsMatch(t, result.WatchedFiles, metadata.WatchedFiles, "Watched files don't match")
		}
	}

	tests.RunTest(test, t, config.AppModule)
}

func TestApplicationMetadataIngestion(t *testing.T) {

	testCases := []struct {
		name     string
		metadata repositories.ApplicationState
	}{
		{
			name: "Ingest metadata with multiple applications",
			metadata: repositories.ApplicationState{
				CollectedAtMs: time.Now().UnixMilli(),
				ClusterId:     "cluster",
				Applications: []repositories.Application{
					{
						Kind: "Deployment",
						Name: "test-dp",
					},
					{
						Kind: "StatefulSet",
						Name: "test-sts",
					},
				},
			},
		},
		{
			name: "Ingest metadata without applications",
			metadata: repositories.ApplicationState{
				CollectedAtMs: time.Now().UnixMilli(),
				ClusterId:     "cluster",
				Applications:  []repositories.Application{},
			},
		},
	}

	type TestDependencies struct {
		fx.In
		Logger                        *zap.Logger
		MetadataService               *services.MetadataService
		ApplicationMetadataBroker     messagebroker.MessageBroker[repositories.ApplicationState]
		ApplicationMetadataRepository *sharedrepo.MongoDbCollection[repositories.ApplicationState]
	}

	test := func(dependencies TestDependencies) {

		for _, test := range testCases {
			dependencies.MetadataService.Init()

			var (
				log      = dependencies.Logger
				metadata = test.metadata
			)
			log.Info("Executing", zap.String("test", test.name))

			dependencies.ApplicationMetadataRepository.DeleteAll()
			dependencies.ApplicationMetadataBroker.Publish("", metadata)

			time.Sleep(15 * time.Second)

			result, err := dependencies.ApplicationMetadataRepository.GetDocument(
				bson.D{
					{Key: "clusterId", Value: metadata.ClusterId},
					{Key: "collectedAtMs", Value: metadata.CollectedAtMs}},
				bson.D{},
			)
			if err != nil {
				log.Error("Error reading application metadata", zap.Error(err))
				t.Error("Error reading application metadata")
			}

			assert.Len(t, result.Applications, len(metadata.Applications), "Invalid number of applications")
			assert.ElementsMatch(t, result.Applications, metadata.Applications, "Applications don't match")
		}
	}

	tests.RunTest(test, t, config.AppModule)
}

func TestNodeMetadataStateUpdate(t *testing.T) {

	testCases := []struct {
		name     string
		metadata []repositories.NodeState
		result   repositories.AggregatedNodeMetadata
	}{
		{
			name: "Update for two nodes with fixed fileset",
			metadata: []repositories.NodeState{
				{
					ClusterId:     "cluster",
					NodeName:      "node1",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2", "file3"},
				},
				{
					ClusterId:     "cluster",
					NodeName:      "node2",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2", "file3"},
				},
			},
			result: repositories.AggregatedNodeMetadata{
				ClusterId: "cluster",
				Metadata: []repositories.NodeMetadata{
					{
						Name:  "node1",
						Files: []interface{}{"file1", "file2", "file3"},
					},
					{
						Name:  "node2",
						Files: []interface{}{"file1", "file2", "file3"},
					},
				},
			},
		},
		{
			name: "Update for a single node with varied fileset",
			metadata: []repositories.NodeState{
				{
					ClusterId:     "cluster",
					NodeName:      "node",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2"},
				},
				{
					ClusterId:     "cluster",
					NodeName:      "node",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2", "file3"},
				},
				{
					ClusterId:     "cluster",
					NodeName:      "node",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file4"},
				},
			},
			result: repositories.AggregatedNodeMetadata{
				ClusterId: "cluster",
				Metadata: []repositories.NodeMetadata{
					{
						Name:  "node",
						Files: []interface{}{"file1", "file2", "file3", "file4"},
					},
				},
			},
		},
	}

	type TestDependencies struct {
		fx.In
		Logger *zap.Logger

		MetadataService *services.MetadataService

		NodeMetadataBroker        messagebroker.MessageBroker[repositories.NodeState]
		NodeMetadataUpdatedBroker messagebroker.MessageBroker[services.NodeMetadataUpdated]

		NodeRepo           *sharedrepo.MongoDbCollection[repositories.NodeState]
		NodeAggregatedRepo *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
	}

	for _, test := range testCases {
		tests.RunTest(func(dependencies TestDependencies) {

			dependencies.MetadataService.Init()

			var (
				metadata = test.metadata
				result   = test.result
				log      = dependencies.Logger
			)

			log.Info("Executing", zap.String("test", test.name))

			dependencies.NodeRepo.DeleteAll()
			dependencies.NodeAggregatedRepo.DeleteAll()

			for _, state := range metadata {
				dependencies.NodeMetadataBroker.Publish("", state)
			}

			msg := make(chan services.NodeMetadataUpdated)
			go dependencies.NodeMetadataUpdatedBroker.Subscribe(context.Background(), msg, make(chan<- error))

			updatedStateEvent := <-msg
			updatedState := updatedStateEvent.Metadata

			assert.Equal(t, result.ClusterId, updatedState.ClusterId, "Invalid clusterId")

			for idx, metadata := range updatedState.Metadata {
				assert.Equal(t, result.Metadata[idx].Name, metadata.Name, "Invalid node name")
				assert.ElementsMatch(t, result.Metadata[idx].Files, metadata.Files, "Invalid node files")
			}

			assert.Len(t, updatedState.Metadata, len(result.Metadata), "Invalid number of nodes")
		},
			t,
			config.AppModule,
		)
	}
}

func TestApplicationMetadataStateUpdate(t *testing.T) {

	testCases := []struct {
		name     string
		metadata []repositories.ApplicationState
		result   repositories.AggregatedApplicationMetadata
	}{
		{
			name: "Application aggregated state with additional applications",
			metadata: []repositories.ApplicationState{
				{
					CollectedAtMs: time.Now().UnixMilli(),
					ClusterId:     "cluster",
					Applications: []repositories.Application{
						{
							Kind: "Deployment",
							Name: "dp",
						},
						{
							Kind: "StatefulSet",
							Name: "sts",
						},
					},
				},
				{
					CollectedAtMs: time.Now().UnixMilli(),
					ClusterId:     "cluster",
					Applications: []repositories.Application{
						{
							Kind: "Deployment",
							Name: "dp",
						},
						{
							Kind: "StatefulSet",
							Name: "sts-2",
						},
						{
							Kind: "DaemonSet",
							Name: "dp-3",
						},
					},
				},
			},
			result: repositories.AggregatedApplicationMetadata{
				ClusterId: "cluster",
				Metadata: []repositories.ApplicationMetadata{
					{
						Kind: "Deployment",
						Name: "dp",
					},
					{
						Kind: "StatefulSet",
						Name: "sts",
					},
					{
						Kind: "StatefulSet",
						Name: "sts-2",
					},
					{
						Kind: "DaemonSet",
						Name: "dp-3",
					},
				},
			},
		},
		{
			name: "Application aggregated state without additional applications",
			metadata: []repositories.ApplicationState{
				{
					CollectedAtMs: time.Now().UnixMilli(),
					ClusterId:     "cluster",
					Applications: []repositories.Application{
						{
							Kind: "Deployment",
							Name: "dp",
						},
						{
							Kind: "StatefulSet",
							Name: "sts",
						},
					},
				},
				{
					CollectedAtMs: time.Now().UnixMilli(),
					ClusterId:     "cluster",
					Applications: []repositories.Application{
						{
							Kind: "Deployment",
							Name: "dp",
						},
						{
							Kind: "StatefulSet",
							Name: "sts",
						},
					},
				},
			},
			result: repositories.AggregatedApplicationMetadata{
				ClusterId: "cluster",
				Metadata: []repositories.ApplicationMetadata{
					{
						Kind: "Deployment",
						Name: "dp",
					},
					{
						Kind: "StatefulSet",
						Name: "sts",
					},
				},
			},
		},
		{
			name: "Application aggregated state with empty applications",
			metadata: []repositories.ApplicationState{
				{
					CollectedAtMs: time.Now().UnixMilli(),
					ClusterId:     "cluster",
					Applications:  []repositories.Application{},
				},
				{
					CollectedAtMs: time.Now().UnixMilli(),
					ClusterId:     "cluster",
					Applications:  []repositories.Application{},
				},
			},
			result: repositories.AggregatedApplicationMetadata{
				ClusterId: "cluster",
				Metadata:  []repositories.ApplicationMetadata{},
			},
		},
	}

	type TestDependencies struct {
		fx.In
		Logger *zap.Logger

		MetadataService *services.MetadataService

		ApplicationMetadataRepository *sharedrepo.MongoDbCollection[repositories.ApplicationState]
		ApplicationAggregatedRepo     *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
		AppRepo                       *sharedrepo.MongoDbCollection[repositories.ApplicationState]

		ApplicationMetadataBroker        messagebroker.MessageBroker[repositories.ApplicationState]
		ApplicationMetadataUpdatedBroker messagebroker.MessageBroker[services.ApplicationMetadataUpdated]
	}

	for _, test := range testCases {
		tests.RunTest(func(dependencies TestDependencies) {

			dependencies.MetadataService.Init()

			var (
				log    = dependencies.Logger
				result = test.result
			)

			log.Info("Executing", zap.String("test", test.name))

			dependencies.AppRepo.DeleteAll()
			dependencies.ApplicationAggregatedRepo.DeleteAll()
			dependencies.ApplicationMetadataRepository.DeleteAll()

			for _, app := range test.metadata {
				dependencies.ApplicationMetadataBroker.Publish("", app)
			}

			msg := make(chan services.ApplicationMetadataUpdated)
			go dependencies.ApplicationMetadataUpdatedBroker.Subscribe(context.Background(), msg, make(chan<- error))

			updatedStateEvent := <-msg

			updatedState := updatedStateEvent.Metadata

			assert.Equal(t, result.ClusterId, updatedState.ClusterId, "Invalid clusterId")
			assert.ElementsMatch(t, result.Metadata, updatedState.Metadata, "Invalid application metadata")
		}, t, config.AppModule)
	}
}

func TestClusterMetadataStateUpdate(t *testing.T) {

	testCases := []struct {
		name     string
		metadata []repositories.NodeState
		result   repositories.AggregatedClusterMetadata
	}{
		{
			name: "test",
			metadata: []repositories.NodeState{
				{
					ClusterId:     "cluster1",
					NodeName:      "node",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2"},
				},
				{
					ClusterId:     "cluster2",
					NodeName:      "node",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2"},
				},
				{
					ClusterId:     "cluster3",
					NodeName:      "node",
					CollectedAtMs: time.Now().UnixMilli(),
					WatchedFiles:  []string{"file1", "file2"},
				},
			},
			result: repositories.AggregatedClusterMetadata{
				Metadata: []repositories.ClusterMetadata{
					{
						ClusterId: "cluster1",
					},
					{
						ClusterId: "cluster2",
					},
					{
						ClusterId: "cluster3",
					},
				},
			},
		},
	}

	type TestDependencies struct {
		fx.In
		MetadataService              *services.MetadataService
		AppRepo                      *sharedrepo.MongoDbCollection[repositories.ApplicationState]
		NodeRepo                     *sharedrepo.MongoDbCollection[repositories.NodeState]
		ClusterAggregatedRepo        *sharedrepo.MongoDbCollection[repositories.AggregatedClusterMetadata]
		NodeMetadataBroker           messagebroker.MessageBroker[repositories.NodeState]
		ClusterMetadataUpdatedBroker messagebroker.MessageBroker[services.ClusterMetadataUpdated]
	}

	for _, test := range testCases {
		tests.RunTest(func(dependencies TestDependencies) {

			dependencies.MetadataService.Init()

			dependencies.ClusterAggregatedRepo.DeleteAll()
			dependencies.NodeRepo.DeleteAll()
			dependencies.AppRepo.DeleteAll()

			for _, state := range test.metadata {
				dependencies.NodeMetadataBroker.Publish("", state)
			}

			msg := make(chan services.ClusterMetadataUpdated)
			go dependencies.ClusterMetadataUpdatedBroker.Subscribe(context.Background(), msg, make(chan<- error))

			updatedState := <-msg

			assert.ElementsMatch(t, test.result.Metadata, updatedState.Metadata.Metadata, "Invalid cluster metadata")
		}, t, config.AppModule)
	}
}
