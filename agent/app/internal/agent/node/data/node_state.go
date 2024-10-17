package data

import "time"

type NodeState struct {
	ClusterName   string   `json:"clusterName"`
	NodeName      string   `json:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles"`
}

func NewNodeState(clusterName, nodeName string, watchedFiles []string) NodeState {
	return NodeState{ClusterName: clusterName, NodeName: nodeName, WatchedFiles: watchedFiles}
}

func (n *NodeState) SetTimestamp() {
	n.CollectedAtMs = time.Now().UnixMilli()
}
