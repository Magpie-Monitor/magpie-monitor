package pods

import "encoding/json"

type PodChunk struct {
	Cluster   string          `json:"cluster"`
	Kind      ApplicationKind `json:"kind"`
	Timestamp int64           `json:"timestamp"`
	Name      string          `json:"name"`
	Namespace string          `json:"namespace"`
	Pods      []Pod           `json:"pods"`
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

type ApplicationKind string

const (
	Deployment  ApplicationKind = "Deployment"
	StatefulSet ApplicationKind = "StatefulSet"
	DaemonSet   ApplicationKind = "DaemonSet"
)

func (o ApplicationKind) String() string {
	switch o {
	case Deployment:
		return "Deployment"
	case StatefulSet:
		return "StatefulSet"
	case DaemonSet:
		return "DaemonSet"
	default:
		return "unknown"
	}
}

func (o ApplicationKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}
