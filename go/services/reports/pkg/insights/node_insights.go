package insights

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/zap"
)

type NodeLogsInsight struct {
	Hostname       string `bson:"hostname"`
	Category       string `bson:"category"`
	Summary        string `bson:"summary"`
	Recommendation string `bson:"recommendation"`
	Source         string `bson:"source"`
	Timestamp      int    `bson:"timestamp"`
}

type NodeInsightsGenerator interface {
	OnDemandNodeInsights(logs []*repositories.NodeLogsDocument) ([]NodeLogsInsight, error)
	ScheduledNodeInsights(
		logs []*repositories.NodeLogsDocument,
		scheduledTime time.Time) ([]NodeLogsInsight, error)
}

type nodeInsightsResponseDto struct {
	Insights []NodeLogsInsight
}

func (g *OpenAiInsightsGenerator) OnDemandNodeInsights(logs []*repositories.NodeLogsDocument) ([]NodeLogsInsight, error) {

	encodedLogs, err := json.Marshal(logs)
	if err != nil {
		g.logger.Error("Failed to encode application logs for application insights", zap.Error(err))
		return nil, err
	}

	openAiResponse, err := g.client.Complete([]*openai.Message{
		{
			Role: "system",
			Content: `You are a kubernetes cluster system administrator. 
			Given a list of logs from nodes which are part of a Kubernetes cluster
			find logs which might suggest any kind of errors or issues. Try to give a possible reason, 
			category of an issue, urgency and possible resolution. Ignore logs which are only 
			informational and are not marked by warnings or errors. Don't provide intruduction. 
			Go straight into describing these logs. As a response for explaination return reports incident, where every inconsitency is an incident. 
			Source is an fragment of a the provided log that you are referencing in summary and recommendation.
			Hostname refers to the Name of the node provided in the log.
			`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`These are logs from my cluster. 
			Please tell me if they might suggest any kind of issues:
			%s`, encodedLogs),
		},
	},
		openai.CreateJsonReponseFormat("node_insights", nodeInsightsResponseDto{}),
	)

	if err != nil {
		g.logger.Error("Failed to get application logs insights from openai client", zap.Error(err))
		return nil, err
	}

	var insights nodeInsightsResponseDto

	err = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &insights)
	if err != nil {
		g.logger.Error("Failed to decode application insights from openai client", zap.Error(err))
		return nil, err
	}

	return insights.Insights, nil
}

func (g *OpenAiInsightsGenerator) ScheduledNodeInsights(logs []*repositories.NodeLogsDocument, scheduledTime time.Time) ([]NodeLogsInsight, error) {
	return nil, nil
}

var _ NodeInsightsGenerator = &OpenAiInsightsGenerator{}
