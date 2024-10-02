package data

import "time"

type NodeState struct {
	NodeName     string   `json:"nodeName"`
	Timestamp    int64    `json:"timestamp"`
	WatchedFiles []string `json:"watchedFiles"`
}

func NewNodeState(nodeName string, watchedFiles []string) NodeState {
	return NodeState{NodeName: nodeName, WatchedFiles: watchedFiles}
}

func (n *NodeState) SetTimestamp() {
	n.Timestamp = time.Now().UnixMicro()
}
