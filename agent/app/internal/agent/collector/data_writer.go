package collector

import (
	"encoding/json"
	nodeData "github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/node/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/agent/pods/data"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/config"
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/remote_write"
	"log"
)

type DataWriter struct {
	podWriter          remote_write.RemoteWriter
	nodeWriter         remote_write.RemoteWriter
	podMetadataWriter  remote_write.RemoteWriter
	nodeMetadataWriter remote_write.RemoteWriter
}

func NewDataWriter(config config.Config) DataWriter {
	return DataWriter{
		podWriter: remote_write.NewStreamWriter(config.Broker.Url, config.Broker.PodTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		nodeWriter: remote_write.NewStreamWriter(config.Broker.Url, config.Broker.NodeTopic, config.Broker.Username,
			config.Broker.Password, config.Broker.BatchSize),
		podMetadataWriter:  remote_write.NewMetadataWriter(config.Global.PodMetadataRemoteWriteUrl),
		nodeMetadataWriter: remote_write.NewMetadataWriter(config.Global.NodeMetadataRemoteWriteUrl),
	}
}

func (d *DataWriter) WriteClusterLogs(chunk data.Chunk) {
	d.writeChunk(chunk, d.podWriter)
}

func (d *DataWriter) WriteClusterMetadata(chunk data.ClusterState) {
	d.writeChunk(chunk, d.podWriter)
}

func (d *DataWriter) WriteNodeLogs(chunk nodeData.Chunk) {
	d.writeChunk(chunk, d.podWriter)
}

func (d *DataWriter) WriteNodeMetadata(chunk nodeData.NodeState) {
	d.writeChunk(chunk, d.podWriter)
}

func (d *DataWriter) writeChunk(chunk interface{}, writer remote_write.RemoteWriter) {
	jsonChunk, err := json.Marshal(chunk)
	if err != nil {
		log.Println("Error converting chunk to JSON: ", err)
	} else {
		writer.Write(string(jsonChunk))
	}
}
