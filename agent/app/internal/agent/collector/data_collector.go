package collector

import (
	"log"

	nodeData "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
)

type DataCollector struct {
	applicationLogsChannel     chan data.Chunk
	applicationMetadataChannel chan data.ClusterState
	nodeLogsChannel            chan nodeData.Chunk
	nodeMetadataChannel        chan nodeData.NodeState
	writer                     DataWriter
}

func NewDataCollector(config *config.Config, c *config.Channels) DataCollector {
	return DataCollector{
		applicationLogsChannel:     c.ApplicationLogsChannel,
		applicationMetadataChannel: c.ApplicationMetadataChannel,
		nodeLogsChannel:            c.NodeLogsChannel,
		nodeMetadataChannel:        c.NodeMetadataChannel,
		writer:                     NewDataWriter(config),
	}
}

func (d *DataCollector) CollectCluster() {
	for {
		select {
		case chunk := <-d.applicationLogsChannel:
			log.Println("Cluster logs collected: ", chunk)
			d.writer.WriteClusterLogs(chunk)
		case chunk := <-d.applicationMetadataChannel:
			log.Println("Cluster metadata collected: ", chunk)
			d.writer.WriteClusterMetadata(chunk)
		}
	}
}

func (d *DataCollector) CollectNodes() {
	for {
		select {
		case chunk := <-d.nodeLogsChannel:
			log.Println("Node logs collected: ", chunk)
			d.writer.WriteNodeLogs(chunk)
		case chunk := <-d.nodeMetadataChannel:
			log.Println("Node metadata collected: ", chunk)
			d.writer.WriteNodeMetadata(chunk)
		}
	}
}
