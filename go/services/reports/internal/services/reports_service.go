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
	Cluster                  string                                          `json:"cluster"`
	FromDate                 int64                                           `json:"fromDate"`
	ToDate                   int64                                           `json:"toDate"`
	ApplicationConfiguration []*repositories.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*repositories.NodeInsightConfiguration        `json:"nodeConfiguration"`
	MaxLength                int                                             `json:"maxLength"`
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

func (s *ReportsService) ScheduleReport(

	ctx context.Context,
	params ReportGenerationFilters,

) (*repositories.Report, error) {
	applicationLogs, err := s.applicationLogsRepository.GetLogs(ctx,
		params.Cluster,
		time.Unix(0, params.FromDate),
		time.Unix(0, params.ToDate))

	nodeLogs, err := s.nodeLogsRepository.GetLogs(ctx,
		params.Cluster,
		time.Unix(0, params.FromDate),
		time.Unix(0, params.ToDate))

	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	filteredApplicationLogs := applicationLogs[0:int(math.Min(float64(params.MaxLength), float64(len(applicationLogs))))]

	applicationInsights, err := s.applicationInsightsGenerator.ScheduleApplicationInsights(
		filteredApplicationLogs,
		params.ApplicationConfiguration, time.Now(),
		params.Cluster,
		params.FromDate,
		params.ToDate,
	)

	filteredNodeLogs := nodeLogs[0:int(math.Min(float64(params.MaxLength), float64(len(applicationLogs))))]

	nodeInsights, err := s.nodeInsightsGenerator.ScheduleNodeInsights(
		filteredNodeLogs,
		params.NodeConfiguration, time.Now(),
		params.Cluster,
		params.FromDate,
		params.ToDate,
	)

	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	report, err := s.reportRepository.InsertReport(ctx, &repositories.Report{
		Cluster:                      params.Cluster,
		Status:                       repositories.ReportState_AwaitingGeneration,
		RequestedAtNs:                time.Now().UnixNano(),
		GeneratedAtNs:                0,
		ScheduledGenerationAtMs:      0,
		FromDateNs:                   params.FromDate,
		ToDateNs:                     params.ToDate,
		ScheduledApplicationInsights: applicationInsights,
		ScheduledNodeInsights:        nodeInsights,
	})

	return report, nil
}

func (s *ReportsService) RetrieveScheduledReport(scheduledReportId string) (*repositories.Report, error) {

	scheduledReport, repoErr := s.reportRepository.GetSingleReport(
		context.TODO(),
		scheduledReportId,
	)

	if repoErr != nil {
		s.logger.Error("Failed to fetch scheduled report", zap.Error(repoErr))
		return nil, repoErr
	}

	scheduledApplicationInsights := scheduledReport.ScheduledApplicationInsights
	insights, err := s.applicationInsightsGenerator.
		GetScheduledApplicationInsights(scheduledApplicationInsights)

	if err != nil {
		s.logger.Error("Failed to get application insights", zap.Error(err))
		return nil, err
	}

	applicationReports, err := s.GetApplicationReportsFromInsights(insights, scheduledApplicationInsights.ApplicationConfiguration)
	if err != nil {
		s.logger.Error("Failed to build applicatin reports from insights", zap.Error(err))
		return nil, err
	}

	scheduledReport.ApplicationReports = applicationReports

	scheduledNodeInsights := scheduledReport.ScheduledNodeInsights
	nodeInsights, err := s.nodeInsightsGenerator.
		GetScheduledNodeInsights(scheduledNodeInsights)

	if err != nil {
		s.logger.Error("Failed to get application insights", zap.Error(err))
		return nil, err
	}

	nodeReports, err := s.GetNodeReportsFromInsights(nodeInsights, scheduledNodeInsights.NodeConfiguration)
	if err != nil {
		s.logger.Error("Failed to build applicatin reports from insights", zap.Error(err))
		return nil, err
	}

	scheduledReport.ApplicationReports = applicationReports
	scheduledReport.NodeReports = nodeReports
	scheduledReport.Status = repositories.ReportState_Generated

	repoErr = s.reportRepository.UpdateReport(context.TODO(), scheduledReport)
	if repoErr != nil {
		s.logger.Error("Failed to update a scheduled report", zap.Error(repoErr))
		return nil, repoErr
	}

	return scheduledReport, nil
}

func (s *ReportsService) GenerateReport(
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

func (s *ReportsService) getApplicationIncidentFromInsight(insight insights.ApplicationInsightsWithMetadata) repositories.ApplicationIncident {

	sources := array.Map(func(metadata insights.ApplicationInsightMetadata) repositories.ApplicationIncidentSource {
		return repositories.ApplicationIncidentSource{
			ContainerName: metadata.ContainerName,
			PodName:       metadata.PodName,
			Content:       metadata.Source,
			Timestamp:     metadata.Timestamp,
		}
	})(insight.Metadata)

	return repositories.ApplicationIncident{
		Category:       insight.Insight.Category,
		Summary:        insight.Insight.Summary,
		Recommendation: insight.Insight.Recommendation,
		Sources:        sources,
	}

}

func (s *ReportsService) getNodeIncidentFromInsight(insight insights.NodeInsightsWithMetadata) repositories.NodeIncident {

	sources := array.Map(func(metadata insights.NodeInsightMetadata) repositories.NodeIncidentSource {
		return repositories.NodeIncidentSource{
			NodeName:  metadata.NodeName,
			Content:   metadata.Source,
			Timestamp: metadata.Timestamp,
		}
	})(insight.Metadata)

	return repositories.NodeIncident{
		Category:       insight.Insight.Category,
		Summary:        insight.Insight.Summary,
		Recommendation: insight.Insight.Recommendation,
		Sources:        sources,
	}

}

func (s *ReportsService) GetApplicationReportsFromInsights(
	applicationInsights []insights.ApplicationInsightsWithMetadata,
	applicationConfiguration []*repositories.ApplicationInsightConfiguration,
) ([]repositories.ApplicationReport, error) {

	insightsByApplication := insights.GroupInsightsByApplication(applicationInsights)

	reports := make([]repositories.ApplicationReport, 0, len(insightsByApplication))
	configByApp := repositories.MapApplicationNameToConfiguration(applicationConfiguration)

	for applicationName, insightsForApplication := range insightsByApplication {

		incidentsFromInsights := array.Map(s.getApplicationIncidentFromInsight)

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

func (s *ReportsService) GetNodeReportsFromInsights(
	nodeInsights []insights.NodeInsightsWithMetadata,
	nodesConfiguration []*repositories.NodeInsightConfiguration,
) ([]repositories.NodeReport, error) {

	insightsByNode := insights.GroupInsightsByNode(nodeInsights)

	reports := make([]repositories.NodeReport, 0, len(insightsByNode))
	configByNode := repositories.MapNodeNameToConfiguration(nodesConfiguration)

	for nodeName, insightsForNode := range insightsByNode {

		nodeIncidentsFromInsights := array.Map(s.getNodeIncidentFromInsight)
		incidents := nodeIncidentsFromInsights(insightsForNode)

		report := repositories.NodeReport{
			Node:      nodeName,
			Incidents: incidents,
		}

		config, ok := configByNode[nodeName]
		if ok {
			report.Precision = config.Precision
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (s *ReportsService) GenerateAndSaveReport(ctx context.Context, params ReportGenerationFilters) (*repositories.Report, error) {

	report, err := s.GenerateReport(ctx,
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

func (s *ReportsService) GenerateApplicationReports(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	applicationConfiguration []*repositories.ApplicationInsightConfiguration,
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

	s.logger.Sugar().Debugf("APPS %v", applicationLogs)

	filteredApplicationLogs := applicationLogs[0:int(math.Min(float64(maxLength), float64(len(applicationLogs))))]

	applicationInsights, err := s.applicationInsightsGenerator.OnDemandApplicationInsights(
		filteredApplicationLogs,
		applicationConfiguration,
	)
	if err != nil {
		s.logger.Error("Failed to generate application insights", zap.Error(err))
		return nil, err
	}
	return s.GetApplicationReportsFromInsights(
		applicationInsights,
		applicationConfiguration,
	)
}

func (s *ReportsService) GenerateNodeReports(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	nodeConfiguration []*repositories.NodeInsightConfiguration,
	maxLength int,
) ([]repositories.NodeReport, error) {

	nodeLogs, err := s.nodeLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)
	if err != nil {
		s.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	filteredNodeLogs := nodeLogs[0:int(math.Min(float64(maxLength), float64(len(nodeLogs))))]

	nodeInsights, err := s.nodeInsightsGenerator.OnDemandNodeInsights(
		filteredNodeLogs,
		nodeConfiguration,
	)
	if err != nil {
		s.logger.Error("Failed to generate application insights", zap.Error(err))
		return nil, err
	}

	s.logger.Sugar().Debugf("NODES %+v", nodeConfiguration[0].CustomPrompt)

	return s.GetNodeReportsFromInsights(
		nodeInsights,
		nodeConfiguration,
	)
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
