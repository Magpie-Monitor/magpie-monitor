package collector

import (
	nodeData "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	"log"
)

type DataCollector struct {
	clusterLogsChannel     chan data.Chunk
	clusterMetadataChannel chan data.ClusterState
	nodeLogsChannel        chan nodeData.Chunk
	nodeMetadataChannel    chan nodeData.NodeState
	writer                 DataWriter
}

func NewDataCollector(config config.Config, channels config.Channels) DataCollector {
	return DataCollector{
		clusterLogsChannel:     channels.ClusterLogsChannel,
		clusterMetadataChannel: channels.ClusterMetadataChannel,
		nodeLogsChannel:        channels.NodeLogsChannel,
		nodeMetadataChannel:    channels.NodeMetadataChannel,
		writer:                 NewDataWriter(config),
	}
}

func (d *DataCollector) CollectCluster() {
	for {
		select {
		case chunk := <-d.clusterLogsChannel:
			log.Println("Cluster logs collected: ", chunk)
			d.writer.WriteClusterLogs(chunk)
		case chunk := <-d.clusterMetadataChannel:
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
