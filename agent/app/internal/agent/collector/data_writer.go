package collector

import (
	"encoding/json"
	"log"

	nodeData "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/broker"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

type DataWriter struct {
	applicationWriter         broker.Broker
	nodeWriter                broker.Broker
	applicationMetadataWriter broker.Broker
	nodeMetadataWriter        broker.Broker
}

func NewDataWriter(config *config.Config) DataWriter {
	return DataWriter{
		applicationWriter: broker.NewStreamWriter(config.Broker.Url, config.Broker.ApplicationTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		nodeWriter: broker.NewStreamWriter(config.Broker.Url, config.Broker.NodeTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		applicationMetadataWriter: broker.NewStreamWriter(config.Broker.Url, config.Broker.ApplicationMetadataTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		nodeMetadataWriter: broker.NewStreamWriter(config.Broker.Url, config.Broker.NodeMetadataTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
	}
}

func (d *DataWriter) WriteApplicationLogs(chunk data.Chunk) {
	d.writeChunk(chunk, d.applicationWriter)
}

func (d *DataWriter) WriteNodeLogs(chunk nodeData.Chunk) {
	d.writeChunk(chunk, d.nodeWriter)
}

func (d *DataWriter) WriteClusterMetadata(chunk data.ApplicationState) {
	d.writeChunk(chunk, d.applicationMetadataWriter)
}

func (d *DataWriter) WriteNodeMetadata(chunk nodeData.NodeState) {
	d.writeChunk(chunk, d.nodeMetadataWriter)
}

func (d *DataWriter) writeChunk(chunk interface{}, writer broker.Broker) {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		log.Println("Error converting chunk to JSON: ", err)
	} else {
		writer.Publish(string(jsonChunk))
	}
}
