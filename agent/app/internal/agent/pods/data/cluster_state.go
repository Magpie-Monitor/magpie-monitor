package data

import (
	v1 "k8s.io/api/apps/v1"
	"time"
)

type ClusterState struct {
	timestamp    int64
	clusterName  string
	applications []Application
}

type Application struct {
	kind ApplicationKind
	name string
}

func NewClusterState(clusterName string) ClusterState {
	return ClusterState{clusterName: clusterName, applications: make([]Application, 0)}
}

type Test interface {
	GetName() string
}

func (c *ClusterState) SetTimestamp() {
	c.timestamp = time.Now().UnixMicro()
}

func (c *ClusterState) AppendDeployments(deployments []v1.Deployment) {
	for _, d := range deployments {
		c.appendApplication(d.Name, Deployment)
	}
}

func (c *ClusterState) AppendStatefulSets(statefulSets []v1.StatefulSet) {
	for _, s := range statefulSets {
		c.appendApplication(s.Name, StatefulSet)
	}
}

func (c *ClusterState) AppendDaemonSets(daemonSets []v1.DaemonSet) {
	for _, d := range daemonSets {
		c.appendApplication(d.Name, DaemonSet)
	}
}

func (c *ClusterState) appendApplication(name string, kind ApplicationKind) {
	app := Application{name: name, kind: kind}
	c.applications = append(c.applications, app)
}
