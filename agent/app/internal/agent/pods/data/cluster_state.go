package data

import (
	"time"

	v1 "k8s.io/api/apps/v1"
)

type ApplicationState struct {
	CollectedAtMs int64         `json:"collectedAtMs"`
	ClusterId     string        `json:"clusterId"`
	Applications  []Application `json:"applications"`
}

type Application struct {
	Kind ApplicationKind `json:"kind"`
	Name string          `json:"name"`
}

func NewApplicationState(clusterId string) ApplicationState {
	return ApplicationState{ClusterId: clusterId, Applications: []Application{}}
}

func (c *ApplicationState) SetTimestamp() {
	c.CollectedAtMs = time.Now().UnixMilli()
}

func (c *ApplicationState) AppendDeployments(deployments *[]v1.Deployment) {
	for _, d := range *deployments {
		c.appendApplication(d.Name, Deployment)
	}
}

func (c *ApplicationState) AppendStatefulSets(statefulSets *[]v1.StatefulSet) {
	for _, s := range *statefulSets {
		c.appendApplication(s.Name, StatefulSet)
	}
}

func (c *ApplicationState) AppendDaemonSets(daemonSets *[]v1.DaemonSet) {
	for _, d := range *daemonSets {
		c.appendApplication(d.Name, DaemonSet)
	}
}

func (c *ApplicationState) appendApplication(name string, kind ApplicationKind) {
	app := Application{Name: name, Kind: kind}
	c.Applications = append(c.Applications, app)
}
