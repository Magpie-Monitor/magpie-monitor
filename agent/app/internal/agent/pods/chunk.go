package pods

type PodChunk struct {
	Cluster   string `json:"cluster"`
	Kind      string `json:"kind"`
	Timestamp int64  `json:"timestamp"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Pods      []Pod  `json:"pods"`
}

type Pod struct {
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
}

type Container struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Content string `json:"content"`
}
