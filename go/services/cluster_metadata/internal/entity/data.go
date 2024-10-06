package entity

type ClusterState struct {
	CollectedAtMs int64         `json:"collectedAtMs"`
	ClusterName   string        `json:"clusterName"`
	Applications  []Application `json:"applications"`
}

type Application struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type NodeState struct {
	NodeName      string   `json:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles"`
}
