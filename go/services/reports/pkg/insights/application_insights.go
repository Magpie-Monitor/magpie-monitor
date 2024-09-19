package insights

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/zap"
)

type ApplicationLogsInsight struct {
	Name           string `bson:"name"`
	Category       string `bson:"category"`
	Summary        string `bson:"summary"`
	Recommendation string `bson:"recommendation"`
	InsightSource  string `bson:"source"`
	// Timestamp      int    `bson:"timestamp"`
}

type ApplicationInsightsGenerator interface {
	OnDemandApplicationInsights(logs []*repositories.ApplicationLogsDocument) ([]ApplicationLogsInsight, error)
	ScheduledApplicationInsights(logs []*repositories.ApplicationLogsDocument, scheduledTime time.Time) ([]ApplicationLogsInsight, error)
}

type applicationInsightsResponseDto struct {
	Insights []ApplicationLogsInsight
}

func (g *OpenAiInsightsGenerator) OnDemandApplicationInsights(logs []*repositories.ApplicationLogsDocument) ([]ApplicationLogsInsight, error) {

	encodedLogs, err := json.Marshal(logs)
	if err != nil {
		g.logger.Error("Failed to encode application logs for application insights", zap.Error(err))
		return nil, err
	}

	openAiResponse, err := g.client.Complete([]*openai.Message{
		{
			Role: "system",
			Content: `You are a kubernetes cluster system administrator. 
			Given a list of logs from a Kubernetes cluster
			find logs which might suggest any kind of errors or issues. Try to give a possible reason, 
			category of an issue, urgency and possible resolution. Ignore logs which are only 
			informational and are not marked by warnings or errors explicitly.  
			Source is an fragment of a the provided log that you are referencing in summary and recommendation. 
			Always declare a unmodified source log with every insight you give. 
			Always give a recommendation on how to resolve the issue. Always give a source. Never repeat insights, ie. 
			if you once use the source do not create an insight for it again. One insight per source. 
			If there are no errors or warnings don't even mention an insight`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`These are logs from my cluster. 
			Please tell me if they might suggest any kind of issues:
			%s`, encodedLogs),
		},
	},
		openai.CreateJsonReponseFormat("application_insights", applicationInsightsResponseDto{}),
	)

	if err != nil {
		g.logger.Error("Failed to get application logs insights from openai client", zap.Error(err))
		return nil, err
	}

	var insights applicationInsightsResponseDto

	err = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &insights)
	if err != nil {
		g.logger.Error("Failed to decode application insights from openai client", zap.Error(err))
		return nil, err
	}

	return insights.Insights, nil
}

func (g *OpenAiInsightsGenerator) ScheduledApplicationInsights(logs []*repositories.ApplicationLogsDocument, scheduledTime time.Time) ([]ApplicationLogsInsight, error) {
	return nil, nil
}

var _ ApplicationInsightsGenerator = &OpenAiInsightsGenerator{}
