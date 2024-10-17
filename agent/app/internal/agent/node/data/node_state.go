package data

import "time"

type NodeState struct {
	ClusterId     string   `json:"clusterId"`
	NodeName      string   `json:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles"`
}

func NewNodeState(clusterId, nodeName string, watchedFiles []string) NodeState {
	return NodeState{ClusterId: clusterId, NodeName: nodeName, WatchedFiles: watchedFiles}
}

func (n *NodeState) SetTimestamp() {
	n.CollectedAtMs = time.Now().UnixMilli()
}
