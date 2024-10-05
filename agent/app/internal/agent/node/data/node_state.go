package data

import "time"

type NodeState struct {
	NodeName      string   `json:"nodeName"`
	CollectedAtMs int64    `json:"collectedAtMs"`
	WatchedFiles  []string `json:"watchedFiles"`
}

func NewNodeState(nodeName string, watchedFiles []string) NodeState {
	return NodeState{NodeName: nodeName, WatchedFiles: watchedFiles}
}

func (n *NodeState) SetTimestamp() {
	n.CollectedAtMs = time.Now().UnixMilli()
}
