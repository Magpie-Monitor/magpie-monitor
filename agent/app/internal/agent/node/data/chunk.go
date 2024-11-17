package data

import "time"

type Chunk struct {
	ClusterId     string `json:"clusterId"`
	Kind          string `json:"kind"`
	CollectedAtMs int64  `json:"collectedAtMs"`
	Name          string `json:"name"`
	Filename      string `json:"filename"`
	Content       string `json:"content"`
}

func NewChunk(clusterId, name, filename, content string) Chunk {
	return Chunk{
		ClusterId:     clusterId,
		Kind:          "Node",
		CollectedAtMs: time.Now().UnixMilli(),
		Name:          name,
		Filename:      filename,
		Content:       content,
	}
}
