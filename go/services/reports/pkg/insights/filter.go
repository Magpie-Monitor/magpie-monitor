package insights

import (
	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"strings"
)

type LogsFilter[A repositories.Log] interface {
	Filter(logs []A) []A
}

type AccuracyFilter[A repositories.Log] struct {
	accuracy         Accuracy
	accuracyKeywords []string
}

func NewAccuracyFilter[T repositories.Log](accuracy Accuracy) *AccuracyFilter[T] {

	return &AccuracyFilter[T]{
		accuracy:         accuracy,
		accuracyKeywords: AccuracyKeywords[accuracy],
	}
}

var AccuracyKeywords = map[Accuracy][]string{
	Accuracy__Low:    LowAccuracyKeywords,
	Accuracy__Medium: MediumAccuracyKeywords,
	Accuracy__High:   HighAccuracyKeywords,
}

var LowAccuracyKeywords = []string{
	"Error", "Exception", "Critical",
	"Warning", "Failed",
	"Authentication", "Authorization",
	"Timeout", "Performance",
	"Security Breach", "Alert", "Audit", "Null",
}

var MediumAccuracyKeywords = append(LowAccuracyKeywords, []string{
	"Info", "User Activity", "Service Call", "Database Query",
	"Reconnect", "Transaction", "Threshold Exceeded", "Data Sync",
	"CPU Usage", "Latency", "Disk I/O", "Health Check", "Backup",
	"Pod Crash", "Node Unreachable", "Container Restart",
	"Service Unavailable", "File Not Found", "Permission Denied",
	"Process Killed", "OOM (Out of Memory)", "Disk Full",
}...)

var HighAccuracyKeywords = append(MediumAccuracyKeywords, []string{
	"Debug", "Notice", "Deprecated", "Retry",
	"Resource Utilization", "Session",
	"Token Expiry", "Credential",
	"Queue Size", "Event Processing", "Memory Leak",
	"Cache Miss", "Cache Hit", "Disk Usage",
	"Connection Lost", "Data Sync", "Transaction",
	"Data Overload", "Memory Usage",
	"Disk I/O", "Throughput", "Data Migration", "Data Ingestion",
	"Network Traffic", "Bandwidth", "Recovery", "Job Queued",
	"Job Completed", "Audit", "Job Failed", "High Latency",
	"Service Unavailable", "Pod Evicted", "Node Disk Pressure",
	"CPU Throttling", "Pod Scheduled", "Ingress Traffic",
	"Egress Traffic", "Node Ready", "Pod Ready", "Pod Pending",
	"Namespace Created", "Namespace Deleted", "ConfigMap Updated",
	"Secret Access", "Persistent Volume Claim Bound",
	"Persistent Volume Claim Released", "DaemonSet Updated",
	"StatefulSet Updated", "ReplicaSet Scaled",
	"Job Completed", "CronJob Triggered",
	"Syslog", "Kernel Panic", "I/O Wait", "Swap Usage",
	"System Load", "Package Installed", "Service Started",
	"User Login", "User Logout", "SSH Access", "Firewall Rule",
	"Mount Point Unavailable", "Network Unreachable",
	"DNS Resolution", "Port Binding", "Kernel Module Loaded",
}...)

func (f *AccuracyFilter[A]) containsKeyword(log *string) bool {

	for _, keyword := range f.accuracyKeywords {
		contains := strings.
			Contains(strings.ToLower(*log), strings.ToLower(keyword))
		if contains == true {
			return true
		}
	}

	return false
}

func (f *AccuracyFilter[A]) Filter(logs []A) []A {
	filtered := make([]A, 0, len(logs)/2)

	for _, log := range logs {
		if f.containsKeyword(log.GetContent()) {
			filtered = append(filtered, log)

		}
	}

	return filtered
}

var _ LogsFilter[*repositories.ApplicationLogsDocument] = &AccuracyFilter[*repositories.ApplicationLogsDocument]{}
