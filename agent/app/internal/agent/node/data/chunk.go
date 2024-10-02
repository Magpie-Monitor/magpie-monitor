package data

type Chunk struct {
	Cluster   string `json:"cluster"`
	Kind      string `json:"kind"`
	Timestamp int64  `json:"timestamp"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Content   string `json:"content"`
}
