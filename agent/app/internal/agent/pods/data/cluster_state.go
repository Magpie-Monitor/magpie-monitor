package data

import (
	"time"

	v1 "k8s.io/api/apps/v1"
)

type ClusterState struct {
	CollectedAtMs int64         `json:"collectedAtMs"`
	ClusterId     string        `json:"clusterId"`
	Applications  []Application `json:"applications"`
}

type Application struct {
	Kind ApplicationKind `json:"kind"`
	Name string          `json:"name"`
}

func NewClusterState(clusterId string) ClusterState {
	return ClusterState{ClusterId: clusterId, Applications: []Application{}}
}

func (c *ClusterState) SetTimestamp() {
	c.CollectedAtMs = time.Now().UnixMilli()
}

func (c *ClusterState) AppendDeployments(deployments *[]v1.Deployment) {
	for _, d := range *deployments {
		c.appendApplication(d.Name, Deployment)
	}
}

func (c *ClusterState) AppendStatefulSets(statefulSets *[]v1.StatefulSet) {
	for _, s := range *statefulSets {
		c.appendApplication(s.Name, StatefulSet)
	}
}

func (c *ClusterState) AppendDaemonSets(daemonSets *[]v1.DaemonSet) {
	for _, d := range *daemonSets {
		c.appendApplication(d.Name, DaemonSet)
	}
}

func (c *ClusterState) appendApplication(name string, kind ApplicationKind) {
	app := Application{Name: name, Kind: kind}
	c.Applications = append(c.Applications, app)
}
