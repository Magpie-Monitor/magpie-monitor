package data

type Chunk struct {
	ClusterId     string `json:"clusterId"`
	Kind          string `json:"kind"`
	CollectedAtMs int64  `json:"collectedAtMs"`
	Name          string `json:"name"`
	Filename      string `json:"filename"`
	Content       string `json:"content"`
}
