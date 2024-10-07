package services

import (
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/IBM/fp-go/array"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

	fromDate := time.Unix(0, params.FromDate)
	toDate := time.Unix(0, params.ToDate)

	applicationLogs, err := s.GetApplicationLogsByParams(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.MaxLength,
	)
	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	nodeLogs, err := s.GetNodeLogsByParams(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.MaxLength,
	)
	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	applicationInsights, err := s.applicationInsightsGenerator.ScheduleApplicationInsights(
		applicationLogs,
		params.ApplicationConfiguration, time.Now(),
		params.Cluster,
		params.FromDate,
		params.ToDate,
	)

	nodeInsights, err := s.nodeInsightsGenerator.ScheduleNodeInsights(
		nodeLogs,
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
		Title:                        s.getTitleForReport(params.Cluster, fromDate, toDate),
		Status:                       repositories.ReportState_AwaitingGeneration,
		RequestedAtNs:                time.Now().UnixNano(),
		ScheduledGenerationAtMs:      time.Now().UnixNano() + time.Hour.Nanoseconds(),
		FromDateNs:                   params.FromDate,
		ToDateNs:                     params.ToDate,
		TotalNodeEntries:             len(nodeLogs),
		TotalApplicationEntries:      len(applicationLogs),
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
	scheduledReport.Urgency = s.getReportUrgencyFromApplicationAndNodeReports(applicationReports, nodeReports)

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

	applicationLogs, err := s.GetApplicationLogsByParams(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.MaxLength,
	)
	if err != nil {
		s.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	applicationReports, err := s.GenerateApplicationReports(
		ctx,
		applicationLogs,
		params.ApplicationConfiguration,
	)
	if err != nil {
		s.logger.Error("Failed to get application reports", zap.Error(err))
		return nil, err
	}

	nodeLogs, err := s.GetNodeLogsByParams(
		ctx,
		params.Cluster,
		fromDate,
		toDate,
		params.MaxLength)
	if err != nil {
		s.logger.Error("Failed to generate node report", zap.Error(err))
		return nil, err
	}

	nodeReports, err := s.GenerateNodeReports(
		ctx,
		nodeLogs,
		params.NodeConfiguration,
	)
	if err != nil {
		s.logger.Error("Failed to get node reports", zap.Error(err))
		return nil, err
	}

	urgency := s.getReportUrgencyFromApplicationAndNodeReports(
		applicationReports,
		nodeReports,
	)

	report := repositories.Report{
		Status:                  repositories.ReportState_Generated,
		Cluster:                 params.Cluster,
		RequestedAtNs:           time.Now().UnixNano(),
		ScheduledGenerationAtMs: time.Now().UnixNano(),
		Title:                   s.getTitleForReport(params.Cluster, fromDate, toDate),
		FromDateNs:              params.FromDate,
		ToDateNs:                params.ToDate,
		NodeReports:             nodeReports,
		ApplicationReports:      applicationReports,
		TotalApplicationEntries: len(applicationLogs),
		TotalNodeEntries:        len(nodeLogs),
		Urgency:                 urgency,
	}

	return &report, nil
}

func (s *ReportsService) getTitleForReport(cluster string, fromDate time.Time, toDate time.Time) string {
	return fmt.Sprintf("Report for %s (%s - %s)",
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
			Image:         metadata.Image,
		}
	})(insight.Metadata)

	return repositories.ApplicationIncident{
		Category:       insight.Insight.Category,
		Summary:        insight.Insight.Summary,
		Recommendation: insight.Insight.Recommendation,
		Urgency:        insight.Insight.Urgency,
		Sources:        sources,
	}

}

// Get maximum of all urgencies from incidents withing passed reports
func (s *ReportsService) getReportUrgencyFromApplicationAndNodeReports(
	applicationReports []repositories.ApplicationReport,
	nodeReports []repositories.NodeReport,
) repositories.Urgency {

	applicationUrgnency := array.Map(func(report repositories.ApplicationReport) []repositories.Urgency {
		return array.Map(func(incident repositories.ApplicationIncident) repositories.Urgency {
			return incident.Urgency
		})(report.Incidents)
	})(applicationReports)

	nodeUrgency := array.Map(func(report repositories.NodeReport) []repositories.Urgency {
		return array.Map(func(incident repositories.NodeIncident) repositories.Urgency {
			return incident.Urgency
		})(report.Incidents)
	})(nodeReports)

	flattenedApplicationUrgency := array.Flatten(applicationUrgnency)
	flattenedNodeUrgency := array.Flatten(nodeUrgency)

	allUrgencies := append(flattenedApplicationUrgency, flattenedNodeUrgency...)

	return slices.Max(allUrgencies)
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
		Urgency:        insight.Insight.Urgency,
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

func (s *ReportsService) GetApplicationLogsByParams(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	maxLength int,
) ([]*sharedrepositories.ApplicationLogsDocument, error) {

	applicationLogs, err := s.applicationLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)
	if err != nil {
		s.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	filteredApplicationLogs := applicationLogs[0:int(math.Min(float64(maxLength), float64(len(applicationLogs))))]

	return filteredApplicationLogs, nil

}

func (s *ReportsService) GenerateApplicationReports(
	ctx context.Context,
	applicationLogs []*sharedrepositories.ApplicationLogsDocument,
	applicationConfiguration []*repositories.ApplicationInsightConfiguration,
) ([]repositories.ApplicationReport, error) {

	applicationInsights, err := s.applicationInsightsGenerator.OnDemandApplicationInsights(
		applicationLogs,
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

func (s *ReportsService) GetNodeLogsByParams(
	ctx context.Context,
	cluster string,
	fromDate time.Time,
	toDate time.Time,
	maxLength int,
) ([]*sharedrepositories.NodeLogsDocument, error) {

	nodeLogs, err := s.nodeLogsRepository.GetLogs(ctx,
		cluster,
		fromDate,
		toDate)
	if err != nil {
		s.logger.Error("Failed to get node logs", zap.Error(err))
		return nil, err
	}

	filteredNodeLogs := nodeLogs[0:int(math.Min(float64(maxLength), float64(len(nodeLogs))))]
	return filteredNodeLogs, nil
}

func (s *ReportsService) GenerateNodeReports(
	ctx context.Context,
	nodeLogs []*sharedrepositories.NodeLogsDocument,
	nodeConfiguration []*repositories.NodeInsightConfiguration,
) ([]repositories.NodeReport, error) {

	nodeInsights, err := s.nodeInsightsGenerator.OnDemandNodeInsights(
		nodeLogs,
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
