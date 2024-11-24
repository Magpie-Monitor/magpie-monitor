package services_test

import (
	"fmt"
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

func TestApplicationMetadataIngestion(t *testing.T) {
	type TestDependencies struct {
		fx.In
		MetadataService               *services.MetadataService
		ApplicationMetadataBroker     *messagebroker.KafkaJsonMessageBroker[repositories.ApplicationState]
		ApplicationMetadataRepository *sharedrepo.MongoDbCollection[repositories.ApplicationState]
	}

	test := func(dependencies TestDependencies) {
		go dependencies.MetadataService.ConsumeApplicationMetadata()

		dependencies.ApplicationMetadataRepository.DeleteAll()

		expectedMetadata := repositories.ApplicationState{
			CollectedAtMs: 1234,
			ClusterId:     "wojtek-test",
			Applications: []repositories.Application{
				repositories.Application{
					Kind: "Deployment",
					Name: "test-dp",
				},
				repositories.Application{
					Kind: "StatefulSet",
					Name: "test-sts",
				},
				repositories.Application{
					Kind: "Deployment",
					Name: "test-dp",
				},
				repositories.Application{
					Kind: "StatefulSet",
					Name: "test-sts",
				},
				repositories.Application{
					Kind: "Deployment",
					Name: "test-dp",
				},
				repositories.Application{
					Kind: "StatefulSet",
					Name: "test-sts",
				},
				repositories.Application{
					Kind: "Deployment",
					Name: "test-dp",
				},
				repositories.Application{
					Kind: "StatefulSet",
					Name: "test-sts",
				},
			},
		}

		dependencies.ApplicationMetadataBroker.Publish("test", expectedMetadata)

		// writer := &kafka.Writer{
		// 	Addr:                   kafka.TCP("kafka:9094"),
		// 	Topic:                  "application_metadata",
		// 	AllowAutoTopicCreation: true,
		// 	Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: "username", Password: "password"}},
		// 	BatchBytes:             int64(5000000),
		// }

		// j, _ := json.Marshal(expectedMetadata)

		// writer.WriteMessages(context.Background(), kafka.Message{Key: []byte("test"), Value: j})

		fmt.Println("test")

		time.Sleep(26 * time.Second)

		metadata, err := dependencies.ApplicationMetadataRepository.GetDocument(bson.D{{Key: "clusterId", Value: "wojtek-test"}, {Key: "collectedAtMs", Value: 1234}}, bson.D{})
		if err != nil {
			t.Fail()
		}

		fmt.Println(metadata)

		docs, _ := dependencies.ApplicationMetadataRepository.GetDocuments(bson.D{}, bson.D{})
		fmt.Println(docs)

		assert.Len(t, metadata.Applications, 8, "Invalid number of applications")
		assert.Equal(t, expectedMetadata.Applications[0].Kind, metadata.Applications[0].Kind)

		dependencies.ApplicationMetadataRepository.DeleteAll()
	}

	tests.RunTest(test, t, config.AppModule)
}

func TestNodeMetadataIngestion(t *testing.T) {
	type TestDependencies struct {
		fx.In
		Logger                 *zap.Logger
		MetadataService        *services.MetadataService
		NodeMetadataBroker     *messagebroker.KafkaJsonMessageBroker[repositories.NodeState]
		NodeMetadataRepository *sharedrepo.MongoDbCollection[repositories.NodeState]
		Creds                  *messagebroker.KafkaCredentials
	}

	test := func(dependencies TestDependencies) {
		go dependencies.MetadataService.ConsumeNodeMetadata()

		_, err := dependencies.NodeMetadataRepository.DeleteAll()
		if err != nil {
			t.Fail()
		}

		expectedMetadata := repositories.NodeState{
			ClusterId:     "wojtek-test",
			NodeName:      "test2",
			CollectedAtMs: 310,
			WatchedFiles: []string{
				"test",
				"test2",
				"test",
				"test2",
				"test",
				"test2",
				"test",
				"test2",
				"test",
				"test2",
			},
		}

		expectedMetadata2 := repositories.NodeState{
			ClusterId:     "wojtek-test",
			NodeName:      "test2",
			CollectedAtMs: 311,
			WatchedFiles: []string{
				"test",
				"test2",
				"test",
				"test2",
				"test",
				"test2",
				"test",
				"test2",
				"test",
				"test2",
			},
		}

		dependencies.NodeMetadataBroker.Publish("test", expectedMetadata)
		dependencies.NodeMetadataBroker.Publish("test", expectedMetadata2)

		// writer := &kafka.Writer{
		// 	Addr:                   kafka.TCP("kafka:9094"),
		// 	Topic:                  "node_metadata",
		// 	AllowAutoTopicCreation: true,
		// 	Transport:              &kafka.Transport{SASL: plain.Mechanism{Username: "username", Password: "password"}},
		// 	BatchBytes:             int64(5000000),
		// }

		// j, _ := json.Marshal(expectedMetadata)

		// writer.WriteMessages(context.Background(), kafka.Message{Key: []byte("test"), Value: j})

		time.Sleep(26 * time.Second)

		metadata, err := dependencies.NodeMetadataRepository.GetDocument(bson.D{{Key: "clusterId", Value: "wojtek-test"}, {Key: "collectedAtMs", Value: 310}}, bson.D{})
		if err != nil {
			t.Fail()
		}

		fmt.Println(metadata, "test")

		assert.Len(t, metadata.WatchedFiles, 10, "Inavlid number of files - TEST")
		assert.Equal(t, expectedMetadata.WatchedFiles[0], metadata.WatchedFiles[0])

		_, err = dependencies.NodeMetadataRepository.DeleteAll()
		if err != nil {
			t.Fail()
		}
	}

	tests.RunTest(test, t, config.AppModule)
}

// func TestApplicationMetadataStateUpdate(t *testing.T) {
// 	type TestDependencies struct {
// 		fx.In
// 		Logger                           *zap.Logger
// 		MetadataService                  *services.MetadataService
// 		ApplicationMetadataBroker        *messagebroker.KafkaJsonMessageBroker[repositories.ApplicationState]
// 		ApplicationMetadataUpdatedBroker *messagebroker.KafkaJsonMessageBroker[services.ApplicationMetadataUpdated]
// 		ApplicationMetadataRepository    *sharedrepo.MongoDbCollection[repositories.ApplicationState]
// 		ApplicationAggregatedRepo        *sharedrepo.MongoDbCollection[repositories.AggregatedApplicationMetadata]
// 		Creds                            *messagebroker.KafkaCredentials
// 	}

// 	test := func(dependencies TestDependencies) {
// 		// c := make(chan services.ApplicationMetadataUpdated)
// 		// err := make(chan error)

// 		// defer close(c)
// 		// defer close(err)

// 		// broker := messagebroker.NewKafkaJsonMessageBroker[services.ApplicationMetadataUpdated](
// 		// 	dependencies.Logger,
// 		// 	dependencies.Creds.Address,
// 		// 	"application_metadata_updated",
// 		// 	dependencies.Creds.Username,
// 		// 	dependencies.Creds.Password,
// 		// )

// 		// reader := kafka.NewReader(
// 		// 	kafka.ReaderConfig{
// 		// 		Brokers:  []string{addr},
// 		// 		Topic:    topic,
// 		// 		MaxBytes: 10e8,
// 		// 		// GroupID:        kafkaBrokerGroupId,
// 		// 		Dialer:         dialer,
// 		// 		CommitInterval: time.Second,
// 		// 	},
// 		// )
// 		// ready := make(chan struct{})
// 		// go func() {

// 		// go dependencies.ApplicationMetadataUpdatedBroker.Subscribe(c, err)
// 		mechanism, err := scram.Mechanism(scram.SHA512, "username", "password")
// 		if err != nil {
// 			panic("Failed to set sasl mechanism for logs ingestion kafka queue")
// 		}

// 		dialer := &kafka.Dialer{
// 			SASLMechanism: mechanism,
// 		}

// 		reader := kafka.NewReader(
// 			kafka.ReaderConfig{
// 				Brokers:        []string{"kafka:9094"},
// 				Topic:          "application_metadata_updated",
// 				MaxBytes:       10e8,
// 				GroupID:        "magpie-monitor",
// 				Dialer:         dialer,
// 				CommitInterval: time.Second,
// 			},
// 		)

// 		_, e := dependencies.ApplicationAggregatedRepo.DeleteAll()
// 		if e != nil {
// 			t.Fail()
// 		}

// 		_, e = dependencies.ApplicationMetadataRepository.DeleteAll()
// 		if e != nil {
// 			t.Fail()
// 		}

// 		expectedMetadata := repositories.ApplicationState{
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			ClusterId:     "test2",
// 			Applications: []repositories.Application{
// 				repositories.Application{
// 					Kind: "Deployment",
// 					Name: "test-dp",
// 				},
// 				repositories.Application{
// 					Kind: "StatefulSet",
// 					Name: "test-sts-1",
// 				},
// 			},
// 		}

// 		expectedMetadata2 := repositories.ApplicationState{
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			ClusterId:     "test2",
// 			Applications: []repositories.Application{
// 				repositories.Application{
// 					Kind: "Deployment",
// 					Name: "test-dp2",
// 				},
// 				repositories.Application{
// 					Kind: "StatefulSet",
// 					Name: "test-sts-1",
// 				},
// 			},
// 		}

// 		expectedMetadata3 := repositories.ApplicationState{
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			ClusterId:     "test2",
// 			Applications: []repositories.Application{
// 				repositories.Application{
// 					Kind: "Deployment",
// 					Name: "test-dp2",
// 				},
// 				repositories.Application{
// 					Kind: "StatefulSet",
// 					Name: "test-sts-1",
// 				},
// 			},
// 		}

// 		dependencies.ApplicationMetadataRepository.InsertDocument(expectedMetadata)
// 		dependencies.ApplicationMetadataRepository.InsertDocument(expectedMetadata2)
// 		dependencies.ApplicationMetadataRepository.InsertDocument(expectedMetadata3)

// 		for {
// 			fmt.Println("here")
// 			msg, _ := reader.ReadMessage(context.Background())
// 			fmt.Println(msg)
// 			break
// 		}

// 		// for {
// 		// 	msg := <-c
// 		// 	fmt.Println(msg)
// 		// 	break
// 		// }

// 		// timeout := time.After(60 * time.Second)
// 		// for {
// 		// 	select {
// 		// 	case metadata := <-c:
// 		// 		fmt.Println(metadata)
// 		// 		return
// 		// 		// case ec := <-err:
// 		// 		// 	fmt.Println(ec)
// 		// 		// 	t.Fail()
// 		// 		// 	return
// 		// 		// case <-timeout:
// 		// 		// 	fmt.Println("timeout")
// 		// 		// 	return
// 		// 	}
// 		// }
// 	}

// 	tests.RunTest(test, t, config.AppModule)
// }

// func TestNodeMetadataStateUpdate(t *testing.T) {

// 	type TestDependencies struct {
// 		fx.In
// 		MetadataService           *services.MetadataService
// 		NodeMetadataBroker        *messagebroker.KafkaJsonMessageBroker[repositories.NodeState]
// 		NodeMetadataUpdatedBroker *messagebroker.KafkaJsonMessageBroker[services.NodeMetadataUpdated]
// 		NodeMetadataRepository    *sharedrepo.MongoDbCollection[repositories.NodeState]
// 		NodeAggregatedRepo        *sharedrepo.MongoDbCollection[repositories.AggregatedNodeMetadata]
// 	}

// 	test := func(dependencies TestDependencies) {
// 		// c := make(chan services.NodeMetadataUpdated)
// 		// err := make(chan error)

// 		// defer close(c)
// 		// defer close(err)

// 		// go dependencies.NodeMetadataUpdatedBroker.Subscribe(c, err)

// 		// ready := make(chan struct{})
// 		// go func() {
// 		// go dependencies.NodeMetadataUpdatedBroker.Subscribe(c, err)
// 		// close(ready)
// 		// }()
// 		// <-ready

// 		// time.Sleep(5 * time.Second)

// 		mechanism, err := scram.Mechanism(scram.SHA512, "username", "password")
// 		if err != nil {
// 			panic("Failed to set sasl mechanism for logs ingestion kafka queue")
// 		}

// 		dialer := &kafka.Dialer{
// 			SASLMechanism: mechanism,
// 		}

// 		reader := kafka.NewReader(
// 			kafka.ReaderConfig{
// 				Brokers:        []string{"kafka:9094"},
// 				Topic:          "node_metadata_updated",
// 				MaxBytes:       10e8,
// 				GroupID:        "magpie-monitor",
// 				Dialer:         dialer,
// 				CommitInterval: time.Second,
// 			},
// 		)

// 		dependencies.NodeAggregatedRepo.DeleteAll()
// 		dependencies.NodeMetadataRepository.DeleteAll()

// 		expectedMetadata := repositories.NodeState{
// 			ClusterId:     "test",
// 			NodeName:      "test2",
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			WatchedFiles:  []string{"test", "test2"},
// 		}

// 		expectedMetadata2 := repositories.NodeState{
// 			ClusterId:     "test",
// 			NodeName:      "test3",
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			WatchedFiles:  []string{"test", "test1", "test5"},
// 		}

// 		dependencies.NodeMetadataRepository.InsertDocument(expectedMetadata)
// 		dependencies.NodeMetadataRepository.InsertDocument(expectedMetadata2)

// 		// go dependencies.NodeMetadataBroker.Publish("", expectedMetadata2)

// 		// msg := <-c
// 		// fmt.Println(msg)

// 		for {
// 			fmt.Println("here")
// 			msg, _ := reader.ReadMessage(context.Background())
// 			fmt.Println(msg)
// 			break
// 		}

// 		// for {
// 		// 	msg := <-c
// 		// 	fmt.Println(msg)
// 		// 	break
// 		// }

// 		// timeout := time.After(60 * time.Second)
// 		// for {
// 		// 	select {
// 		// 	case metadata := <-c:
// 		// 		fmt.Println(metadata)
// 		// 		return
// 		// 		// case ec := <-err:
// 		// 		// 	fmt.Println(ec)
// 		// 		// 	t.Fail()
// 		// 		// 	return
// 		// 		// case <-timeout:
// 		// 		// 	fmt.Println("timeout")
// 		// 		// 	return
// 		// 	}
// 		// }

// 		// 	case <-time.After(100 * time.Second):
// 		// 		fmt.Println("timeout")
// 		// 		t.Fail()
// 		// 	}
// 		// }

// 		// dependencies.NodeMetadataRepository.DeleteAll()
// 	}

// 	tests.RunTest(test, t, config.AppModule)
// }

// func TestClusterMetadataStateUpdate(t *testing.T) {
// 	type TestDependencies struct {
// 		fx.In
// 		MetadataService              *services.MetadataService
// 		NodeMetadataBroker           *messagebroker.KafkaJsonMessageBroker[repositories.NodeState]
// 		ClusterMetadataUpdatedBroker *messagebroker.KafkaJsonMessageBroker[services.ClusterMetadataUpdated]
// 	}

// 	test := func(dependencies TestDependencies) {
// 		expectedMetadata := repositories.NodeState{
// 			ClusterId:     "test",
// 			NodeName:      "test2",
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			WatchedFiles:  []string{"test", "test2"},
// 		}

// 		expectedMetadata2 := repositories.NodeState{
// 			ClusterId:     "test22",
// 			NodeName:      "test2",
// 			CollectedAtMs: time.Now().UnixMilli(),
// 			WatchedFiles:  []string{"test", "test2"},
// 		}

// 		go dependencies.MetadataService.ConsumeNodeMetadata()
// 		go dependencies.MetadataService.PollForClusterStateChange()

// 		dependencies.NodeMetadataBroker.Publish("", expectedMetadata)
// 		dependencies.NodeMetadataBroker.Publish("", expectedMetadata2)

// 		time.Sleep(10 * time.Second)

// 		c := make(chan services.ClusterMetadataUpdated)
// 		err := make(chan error)

// 		go dependencies.ClusterMetadataUpdatedBroker.Subscribe(c, err)

// 		select {
// 		case metadata := <-c:
// 			fmt.Println(metadata)
// 		case err := <-err:
// 			fmt.Println(err)
// 			t.Fail()
// 		}

// 	}

// 	tests.RunTest(test, t, config.AppModule)
// }
