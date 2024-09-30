package services

import (
	"context"
	"fmt"
	"github.com/IBM/fp-go/array"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"math"
	"time"
)

type ReportGenerationFilters struct {
	Cluster                  string                                      `json:"cluster"`
	FromDate                 int64                                       `json:"fromDate"`
	ToDate                   int64                                       `json:"toDate"`
	ApplicationConfiguration []*insights.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*insights.NodeInsightConfiguration        `json:"nodeConfiguration"`
	MaxLength                int                                         `json:"maxLength"`
}

type ReportsService struct {
	logger                       *zap.Logger
	reportRepository             repositories.ReportRepository
	applicationLogsRepository    sharedrepositories.ApplicationLogsRepository
	nodeLogsRepository           sharedrepositories.NodeLogsRepository
	applicationInsightsGenerator insights.ApplicationInsightsGenerator
	nodeInsightsGenerator        insights.NodeInsightsGenerator
}

type ReportsServerParams struct {
	fx.In
	Logger                       *zap.Logger
	ReportRepository             repositories.ReportRepository
	ApplicationLogsRepository    sharedrepositories.ApplicationLogsRepository
	NodeLogsRepository           sharedrepositories.NodeLogsRepository
	ApplicationInsightsGenerator insights.ApplicationInsightsGenerator
	NodeInsightsGenerator        insights.NodeInsightsGenerator
}

func NewReportsService(p ReportsServerParams) *ReportsService {
	return &ReportsService{
		logger:                       p.Logger,
		reportRepository:             p.ReportRepository,
		applicationLogsRepository:    p.ApplicationLogsRepository,
		nodeLogsRepository:           p.NodeLogsRepository,
		applicationInsightsGenerator: p.ApplicationInsightsGenerator,
		nodeInsightsGenerator:        p.NodeInsightsGenerator,
	}
}

func (s *ReportsService) GenerateCompleteOnDemandReport(
	ctx context.Context,
	params ReportGenerationFilters) (*repositories.Report, error) {

	fromDate := time.Unix(0, params.FromDate)
	toDate := time.Unix(0, params.ToDate)

	applicationReports, err := s.GenerateApplicationReports(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.ApplicationConfiguration,
		params.MaxLength,
	)
	if err != nil {
		s.logger.Error("Failed to get application reports", zap.Error(err))
		return nil, err
	}

	nodeReports, err := s.GenerateNodeReports(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.NodeConfiguration,
		params.MaxLength,
	)
	if err != nil {
		s.logger.Error("Failed to get node reports", zap.Error(err))
		return nil, err
	}

	report := repositories.Report{
		Status:        repositories.ReportState_Generated,
		Cluster:       params.Cluster,
		RequestedAtNs: time.Now().UnixNano(),
		//TODO: Generate the report entity first and assign the insights
		//later once the Batch API is implemented.
		GeneratedAtNs:           time.Now().UnixNano(),
		ScheduledGenerationAtMs: time.Now().UnixNano(),
		Title:                   s.getTitleForReport(params.Cluster, fromDate, toDate),
		FromDateNs:              params.FromDate,
		ToDateNs:                params.ToDate,
		NodeReports:             nodeReports,
		ApplicationReports:      applicationReports,
	}

	return &report, nil
}

func (s *ReportsService) getTitleForReport(cluster string, fromDate time.Time, toDate time.Time) string {
	return fmt.Sprintf("On-Demand report for %s (%s - %s)",
		cluster,
		fmt.Sprintf("%d.%d", fromDate.Month(), fromDate.Year()),
		fmt.Sprintf("%d.%d", toDate.Month(), toDate.Year()),
	)
}

func (s *ReportsService) GenerateApplicationReports(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	applicationConfiguration []*insights.ApplicationInsightConfiguration,
	maxLength int,
) ([]repositories.ApplicationReport, error) {

	applicationLogs, err := s.applicationLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)

	if err != nil {
		s.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	filteredApplicationLogs := applicationLogs[0:int(math.Min(float64(maxLength), float64(len(applicationLogs))))]

	applicationInsights, err := s.applicationInsightsGenerator.OnDemandApplicationInsights(
		filteredApplicationLogs,
		applicationConfiguration,
	)
	if err != nil {
		s.logger.Error("Failed to generate application insights", zap.Error(err))
		return nil, err
	}

	insightsByApplication := make(map[string][]insights.ApplicationInsightsWithMetadata)

	for _, insight := range applicationInsights {
		applicationName := insight.Metadata.ApplicationName
		insightsByApplication[applicationName] = append(insightsByApplication[applicationName], insight)
	}

	reports := make([]repositories.ApplicationReport, 0, len(insightsByApplication))

	configByApp := insights.MapApplicationNameToConfiguration(applicationConfiguration)

	for applicationName, insightsForApplication := range insightsByApplication {

		incidentsFromInsights := array.Map(func(insight insights.ApplicationInsightsWithMetadata) repositories.ApplicationIncident {
			return repositories.ApplicationIncident{
				Category:       insight.Insight.Category,
				Summary:        insight.Insight.Summary,
				Recommendation: insight.Insight.Recommendation,
				Timestamp:      insight.Metadata.Timestamp,
				PodName:        insight.Metadata.PodName,
				ContainerName:  insight.Metadata.ContainerName,
				Source:         insight.Metadata.Source,
			}
		})

		incidents := incidentsFromInsights(insightsForApplication)
		report := repositories.ApplicationReport{
			ApplicationName: applicationName,
			Incidents:       incidents,
		}

		config, ok := configByApp[applicationName]
		if ok {
			report.Precision = config.Precision
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil

}

func (s *ReportsService) GenerateAndSaveReport(ctx context.Context, params ReportGenerationFilters) (*repositories.Report, error) {

	report, err := s.GenerateCompleteOnDemandReport(ctx,
		ReportGenerationFilters{
			Cluster:                  params.Cluster,
			FromDate:                 params.FromDate,
			ToDate:                   params.ToDate,
			MaxLength:                params.MaxLength,
			ApplicationConfiguration: params.ApplicationConfiguration,
			NodeConfiguration:        params.NodeConfiguration,
		})
	if err != nil {
		s.logger.Error("Failed to generate report", zap.Error(err))
		return nil, err
	}

	insertedReport, repositoryErr := s.InsertReport(ctx, report)
	if repositoryErr != nil {
		s.logger.Error("Failed to save a report", zap.Error(err))
		return nil, err
	}

	return insertedReport, err
}

func (s *ReportsService) GenerateNodeReports(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	nodeConfiguration []*insights.NodeInsightConfiguration,
	maxLength int,
) ([]repositories.NodeReport, error) {

	nodeLogs, err := s.nodeLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)

	if err != nil {
		s.logger.Error("Failed to get node logs", zap.Error(err))
		return nil, err
	}

	filteredNodeLogs := nodeLogs[0:int(math.Min(float64(maxLength), float64(len(nodeLogs))))]

	nodeInsights, err := s.nodeInsightsGenerator.OnDemandNodeInsights(
		filteredNodeLogs,
		nodeConfiguration,
	)
	if err != nil {
		s.logger.Error("Failed to generate node insights", zap.Error(err))
		return nil, err
	}

	insightsByHostname := make(map[string][]insights.NodeInsightsWithMetadata)

	for _, insight := range nodeInsights {
		hostname := insight.Metadata.NodeName
		insightsByHostname[hostname] = append(insightsByHostname[hostname], insight)
	}

	reports := make([]repositories.NodeReport, 0, len(insightsByHostname))

	configByApp := insights.MapNodeNameToConfiguration(nodeConfiguration)

	for hostname, insightsForNode := range insightsByHostname {

		incidentsFromInsights := array.Map(func(insight insights.NodeInsightsWithMetadata) repositories.NodeIncident {
			return repositories.NodeIncident{
				Category:       insight.Insight.Category,
				Summary:        insight.Insight.Summary,
				Recommendation: insight.Insight.Recommendation,
				Timestamp:      insight.Metadata.Timestamp,
				Source:         insight.Metadata.Source,
			}
		})

		incidents := incidentsFromInsights(insightsForNode)
		report := repositories.NodeReport{
			Host:      hostname,
			Incidents: incidents,
		}

		config, ok := configByApp[hostname]
		if ok {
			report.Precision = config.Precision
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil

}

func (s *ReportsService) GetSingleReport(ctx context.Context, id string) (*repositories.Report, *repositories.ReportRepositoryError) {
	return s.reportRepository.GetSingleReport(ctx, id)
}

func (s *ReportsService) GetAllReports(ctx context.Context, filter repositories.FilterParams) ([]*repositories.Report, *repositories.ReportRepositoryError) {
	return s.reportRepository.GetAllReports(ctx, filter)
}

func (s *ReportsService) InsertReport(ctx context.Context, report *repositories.Report) (*repositories.Report, *repositories.ReportRepositoryError) {
	return s.reportRepository.InsertReport(ctx, report)
}
