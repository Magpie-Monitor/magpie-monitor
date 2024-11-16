package services

import (
	"context"
	"fmt"
	"github.com/IBM/fp-go/array"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	sharedrepositories "github.com/Magpie-Monitor/magpie-monitor/pkg/repositories"
	incidentcorrelation "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/incident_correlation"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/insights"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

const (
	REPORTS_POLLING_INTERVAL_SECONDS_KEY = "REPORTS_POLLING_INTERVAL_SECONDS"
)

type ReportGenerationFilters struct {
	ClusterId                string                                      `json:"clusterId"`
	CorrelationId            string                                      `json:"correlationId"`
	SinceMs                  int64                                       `json:"sinceMs"`
	ToMs                     int64                                       `json:"toMs"`
	ApplicationConfiguration []*insights.ApplicationInsightConfiguration `json:"applicationConfiguration"`
	NodeConfiguration        []*insights.NodeInsightConfiguration        `json:"nodeConfiguration"`
}

type ReportsService struct {
	logger                       *zap.Logger
	reportRepository             repositories.ReportRepository
	applicationLogsRepository    sharedrepositories.ApplicationLogsRepository
	nodeLogsRepository           sharedrepositories.NodeLogsRepository
	applicationInsightsGenerator insights.ApplicationInsightsGenerator
	nodeInsightsGenerator        insights.NodeInsightsGenerator
	incidentMerger               incidentcorrelation.IncidentMerger
	pollingIntervalSeconds       int
}

type ReportsServerParams struct {
	fx.In
	Logger                       *zap.Logger
	ReportRepository             repositories.ReportRepository
	ApplicationLogsRepository    sharedrepositories.ApplicationLogsRepository
	NodeLogsRepository           sharedrepositories.NodeLogsRepository
	ApplicationInsightsGenerator insights.ApplicationInsightsGenerator
	IncidentMerger               incidentcorrelation.IncidentMerger
	NodeInsightsGenerator        insights.NodeInsightsGenerator
}

func NewReportsService(p ReportsServerParams) *ReportsService {

	envs.ValidateEnvs("Failed to create ReportsService, missing envs", []string{
		REPORTS_POLLING_INTERVAL_SECONDS_KEY,
	})

	reportsPollingIntervalSeconds := os.Getenv(REPORTS_POLLING_INTERVAL_SECONDS_KEY)
	pollingIntervalSecondsInt, err := strconv.Atoi(reportsPollingIntervalSeconds)
	if err != nil {
		panic(fmt.Sprintf("%s is not a number", REPORTS_POLLING_INTERVAL_SECONDS_KEY))
	}

	return &ReportsService{
		logger:                       p.Logger,
		reportRepository:             p.ReportRepository,
		applicationLogsRepository:    p.ApplicationLogsRepository,
		nodeLogsRepository:           p.NodeLogsRepository,
		applicationInsightsGenerator: p.ApplicationInsightsGenerator,
		nodeInsightsGenerator:        p.NodeInsightsGenerator,
		pollingIntervalSeconds:       pollingIntervalSecondsInt,
		incidentMerger:               p.IncidentMerger,
	}
}

func (s *ReportsService) ScheduleReport(
	ctx context.Context,
	params ReportGenerationFilters,
) (*repositories.Report, error) {

	sinceDate := time.UnixMilli(params.SinceMs)
	toDate := time.UnixMilli(params.ToMs)

	applicationLogs, err := s.GetApplicationLogsByParams(
		ctx,
		params.ClusterId,
		sinceDate,
		toDate,
	)
	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	nodeLogs, err := s.GetNodeLogsByParams(
		ctx,
		params.ClusterId,
		sinceDate,
		toDate)
	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	var aggregatedApplicationInsights = insights.ScheduledApplicationInsights{}

	for {
		if !applicationLogs.HasNextBatch() {
			break
		}

		nextLogsBatch, err := applicationLogs.GetNextBatch()
		if err != nil {
			s.logger.Error("Failed to get next batch of logs")
		}

		applicationInsights, err := s.applicationInsightsGenerator.ScheduleApplicationInsights(
			nextLogsBatch,
			params.ApplicationConfiguration, time.Now(),
			params.ClusterId,
			params.SinceMs,
			params.ToMs,
		)

		aggregatedApplicationInsights = insights.ScheduledApplicationInsights{
			ScheduledJobIds:          append(aggregatedApplicationInsights.ScheduledJobIds, applicationInsights.ScheduledJobIds...),
			SinceMs:                  applicationInsights.SinceMs,
			ToMs:                     applicationInsights.ToMs,
			ClusterId:                applicationInsights.ClusterId,
			ApplicationConfiguration: applicationInsights.ApplicationConfiguration,
		}
	}

	var aggregatedNodeInsights = insights.ScheduledNodeInsights{}

	for {
		if !applicationLogs.HasNextBatch() {
			break
		}

		nextLogsBatch, err := nodeLogs.GetNextBatch()
		if err != nil {
			s.logger.Error("Failed to get next batch of logs")
		}

		nodeInsights, err := s.nodeInsightsGenerator.ScheduleNodeInsights(
			nextLogsBatch,
			params.NodeConfiguration, time.Now(),
			params.ClusterId,
			params.SinceMs,
			params.ToMs,
		)

		aggregatedNodeInsights = insights.ScheduledNodeInsights{
			ScheduledJobIds:   append(aggregatedNodeInsights.ScheduledJobIds, nodeInsights.ScheduledJobIds...),
			SinceMs:           nodeInsights.SinceMs,
			ToMs:              nodeInsights.ToMs,
			ClusterId:         nodeInsights.ClusterId,
			NodeConfiguration: nodeInsights.NodeConfiguration,
		}
	}

	if err != nil {
		s.logger.Error("Failed to fetch application logs", zap.Error(err))
		return nil, err
	}

	report, err := s.reportRepository.InsertReport(ctx, &repositories.Report{
		ClusterId:               params.ClusterId,
		CorrelationId:           params.CorrelationId,
		Title:                   s.getTitleForReport(params.ClusterId, sinceDate, toDate),
		Status:                  repositories.ReportState_AwaitingGeneration,
		RequestedAtMs:           time.Now().UnixMilli(),
		ScheduledGenerationAtMs: time.Now().UnixMilli() + time.Hour.Milliseconds(),
		SinceMs:                 params.SinceMs,
		ToMs:                    params.ToMs,
		// TotalNodeEntries:             len(aggregatedNodeInsights),
		// TotalApplicationEntries:      len(applicationLogs),
		ScheduledApplicationInsights: &aggregatedApplicationInsights,
		ScheduledNodeInsights:        &aggregatedNodeInsights,
		NodeReports:                  []*repositories.NodeReport{},
		ApplicationReports:           []*repositories.ApplicationReport{},
	})

	return report, nil
}

type ReportGenerationError struct {
	Msg    string
	Report *repositories.Report
}

func (s *ReportsService) PollReportsPendingIncidentMerge(ctx context.Context, reportsChn chan<- *repositories.Report, errChn chan<- *ReportGenerationError) {

	for {
		reports, err := s.reportRepository.GetPendingIncidentMergingReports(ctx)

		if err != nil {
			s.logger.Error("Failed to get reports pending incident merging", zap.Error(err))
			continue
		}

		for _, report := range reports {

			applicationJobs := report.ScheduledApplicationIncidentMergerJobs
			nodeJobs := report.ScheduledNodeIncidentMergerJobs

			applicationJobsFinished, err := s.incidentMerger.AreAllJobsFinished(applicationJobs)
			if err != nil {
				s.logger.Error("Failed to check if application incident merger jobs were finished", zap.Error(err))
				continue
			}

			nodeJobsFinished, err := s.incidentMerger.AreAllJobsFinished(nodeJobs)
			if err != nil {
				s.logger.Error("Failed to check if application incident merger jobs were finished")
				continue
			}

			if !(nodeJobsFinished && applicationJobsFinished) {
				continue
			}

			err = s.MergeApplicationIncidents(applicationJobs, report)
			if err != nil {
				s.logger.Error("Failed to merge application incidents", zap.Error(err))
				continue
			}

			err = s.MergeNodeIncidents(nodeJobs, report)
			if err != nil {
				s.logger.Error("Failed to merge node incidents", zap.Error(err))
				continue
			}

			report.Status = repositories.ReportState_Generated
			repoErr := s.reportRepository.UpdateReport(context.TODO(), report)
			if repoErr != nil {
				s.logger.Error("Failed to update report after incident merge")
			}

			s.logger.Info("Merged incidents of a report", zap.Any("reportId", report.Id))
			reportsChn <- report
		}

		time.Sleep(time.Second * time.Duration(s.pollingIntervalSeconds))
	}
}

func (s *ReportsService) MergeApplicationIncidents(applicationMergerJobs []*repositories.ScheduledIncidentMergerJob, report *repositories.Report) error {

	for _, job := range applicationMergerJobs {

		applicationMergerGroups, err := s.incidentMerger.TryGettingIncidentMergerJobIfFinished(job)
		if err != nil {
			s.logger.Error("Failed to get application merger group", zap.Error(err))
			continue
		}

		for idx, applicationReport := range report.ApplicationReports {
			mergerGroups, ok := applicationMergerGroups[applicationReport.ApplicationName]
			if !ok {
				s.logger.Info("Application is not present in merger groups, skipping")
				continue
			}

			newApplicationIncidents := incidentcorrelation.MergeApplicationIncidentsByGroups(mergerGroups, applicationReport)
			insertedIncidents, err := s.reportRepository.InsertApplicationIncidents(context.TODO(), newApplicationIncidents)
			if err != nil {
				s.logger.Error("Failed to insert merged application incidents", zap.Error(err))
				return err
			}
			report.ApplicationReports[idx].Incidents = insertedIncidents

		}
	}

	return nil
}

func (s *ReportsService) MergeNodeIncidents(nodeIncidentMergerJobs []*repositories.ScheduledIncidentMergerJob, report *repositories.Report) error {

	for _, job := range nodeIncidentMergerJobs {

		nodeMergerGroups, err := s.incidentMerger.TryGettingIncidentMergerJobIfFinished(job)
		if err != nil {
			s.logger.Error("Failed to get application merger group", zap.Error(err))
		}

		for idx, nodeReport := range report.NodeReports {
			mergerGroups, ok := nodeMergerGroups[nodeReport.Node]
			if !ok {
				s.logger.Info("Application is not present in merger groups, skipping")
				continue
			}

			newNodeIncidents := incidentcorrelation.MergeNodeIncidentsByGroups(mergerGroups, nodeReport)
			insertedIncidents, err := s.reportRepository.InsertNodeIncidents(context.TODO(), newNodeIncidents)
			if err != nil {
				s.logger.Error("Failed to insert merged application incidents", zap.Error(err))
				return err
			}

			report.NodeReports[idx].Incidents = insertedIncidents
		}
	}

	return nil
}

func (s *ReportsService) PollReportsPendingGeneration(ctx context.Context, reportsChn chan<- *repositories.Report, errChn chan<- *ReportGenerationError) {

	pendingReports := make(map[string]bool)

	for {
		reports, err := s.reportRepository.GetPendingGenerationReports(ctx)

		if err != nil {
			s.logger.Error("Failed to get pending reports", zap.Error(err))
			continue
		}

		for _, report := range reports {
			if isAlreadyPending := pendingReports[report.Id]; isAlreadyPending == true {
				continue
			}

			pendingReports[report.Id] = true
			go func() {

				s.logger.Info("Pending reports", zap.Any("id", report.Id),
					zap.Any("status", report.Status),
					zap.Any("scheduled", report.ScheduledApplicationInsights),
				)

				applicationInsights, err := s.applicationInsightsGenerator.AwaitScheduledApplicationInsights(report.ScheduledApplicationInsights)
				if err != nil {
					s.logger.Error("Failed to await for application insights", zap.Error(err), zap.Any("insights", report.ScheduledApplicationInsights))
					errChn <- &ReportGenerationError{
						Msg:    fmt.Sprintf("Failed to await for application insights: %s", err.Error()),
						Report: report,
					}
					return
				}

				nodeInsights, err := s.nodeInsightsGenerator.AwaitScheduledNodeInsights(report.ScheduledNodeInsights)
				if err != nil {
					s.logger.Error("Failed to await for node insights", zap.Error(err), zap.Any("insights", report.ScheduledNodeInsights))

					errChn <- &ReportGenerationError{
						Msg:    fmt.Sprintf("Failed to await for application insights: %s", err.Error()),
						Report: report,
					}
					return
				}

				report, err := s.CompletePendingReport(report.Id, applicationInsights, nodeInsights)
				if err != nil {
					s.logger.Error("Failed to complete pending report", zap.Error(err), zap.Any("insights", report.ScheduledNodeInsights))
					errChn <- &ReportGenerationError{
						Msg:    fmt.Sprintf("Failed to complete pending report: %s", err.Error()),
						Report: report,
					}
					return
				}

				delete(pendingReports, report.Id)

				s.logger.Info("Completed report", zap.Any("reportId", report.Id))
				reportsChn <- report
			}()
		}

		time.Sleep(time.Second * time.Duration(s.pollingIntervalSeconds))
	}
}

func (s *ReportsService) InsertNodeIncidents(reports []*repositories.NodeReport) error {
	for idx, report := range reports {
		insertedIncidents, err := s.reportRepository.InsertNodeIncidents(context.TODO(), report.Incidents)
		if err != nil {
			s.logger.Error("Failed to insert node incidents", zap.Error(err))
			return err
		}

		reports[idx].Incidents = insertedIncidents
	}

	return nil
}

func (s *ReportsService) InsertApplicationIncidents(reports []*repositories.ApplicationReport) error {
	for idx, report := range reports {
		insertedIncidents, err := s.reportRepository.InsertApplicationIncidents(context.TODO(), report.Incidents)
		if err != nil {
			s.logger.Error("Failed to insert node incidents", zap.Error(err))
			return err
		}

		reports[idx].Incidents = insertedIncidents
	}

	return nil
}

func (s *ReportsService) CompletePendingReport(reportId string, applicationInsights []insights.ApplicationInsightsWithMetadata,
	nodeInsights []insights.NodeInsightsWithMetadata) (*repositories.Report, error) {

	scheduledReport, repoErr := s.reportRepository.GetSingleReport(
		context.TODO(),
		reportId,
	)
	if repoErr != nil {
		s.logger.Error("Failed to fetch scheduled report", zap.Error(repoErr))
		return nil, repoErr
	}

	applicationReports, err := s.GetApplicationReportsFromInsights(applicationInsights, scheduledReport.ScheduledApplicationInsights.ApplicationConfiguration)
	if err != nil {
		s.logger.Error("Failed to build applicatin reports from insights", zap.Error(err))
		return nil, err
	}

	nodeReports, err := s.GetNodeReportsFromInsights(nodeInsights, scheduledReport.ScheduledNodeInsights.NodeConfiguration)
	if err != nil {
		s.logger.Error("Failed to build applicatin reports from insights", zap.Error(err))
		return nil, err
	}

	err = s.InsertApplicationIncidents(applicationReports)
	if err != nil {
		s.logger.Error("Failed to insert node incidents")
		return nil, err
	}

	err = s.InsertNodeIncidents(nodeReports)
	if err != nil {
		s.logger.Error("Failed to insert application incidents")
		return nil, err
	}

	scheduledApplicationIncidentsMergerJobs, err := s.ScheduleApplicationIncidentsMerge(applicationReports)
	if err != nil {
		s.logger.Error("Failed to schedule application incident merge")
		return nil, err
	}

	scheduledNodeIncidentsMergerJobs, err := s.ScheduleNodeIncidentsMerge(nodeReports)
	if err != nil {
		s.logger.Error("Failed to schedule node incident merge")
		return nil, err
	}

	scheduledReport.ApplicationReports = applicationReports
	scheduledReport.NodeReports = nodeReports
	scheduledReport.Status = repositories.ReportState_AwaitingIncidentMerging
	scheduledReport.Urgency = s.getReportUrgencyFromApplicationAndNodeReports(applicationReports, nodeReports)
	scheduledReport.AnalyzedNodes = len(scheduledReport.ScheduledNodeInsights.NodeConfiguration)
	scheduledReport.AnalyzedApplications = len(scheduledReport.ScheduledApplicationInsights.ApplicationConfiguration)
	scheduledReport.ScheduledNodeIncidentMergerJobs = scheduledNodeIncidentsMergerJobs
	scheduledReport.ScheduledApplicationIncidentMergerJobs = scheduledApplicationIncidentsMergerJobs

	repoErr = s.reportRepository.UpdateReport(context.TODO(), scheduledReport)
	if repoErr != nil {
		s.logger.Error("Failed to update a scheduled report", zap.Error(repoErr))
		return nil, repoErr
	}

	updatedReport, repoErr := s.reportRepository.GetSingleReport(context.TODO(), scheduledReport.Id)
	if repoErr != nil {
		s.logger.Error("Failed to update a scheduled report", zap.Error(repoErr))
		return nil, repoErr
	}

	return updatedReport, nil
}

func (s *ReportsService) ScheduleApplicationIncidentsMerge(
	applicationReports []*repositories.ApplicationReport) ([]*repositories.ScheduledIncidentMergerJob, error) {

	incidentsByApps := make(map[string][]repositories.Incident)
	for _, report := range applicationReports {
		incidentsByApps[report.ApplicationName] = incidentcorrelation.ConvertConcreteIncidentArrayIntoIncidents(report.Incidents)
	}

	return s.incidentMerger.ScheduleIncidentsMerge(incidentsByApps)
}

func (s *ReportsService) ScheduleNodeIncidentsMerge(nodeReports []*repositories.NodeReport) ([]*repositories.ScheduledIncidentMergerJob, error) {

	incidentsByNodes := make(map[string][]repositories.Incident)
	for _, report := range nodeReports {
		incidentsByNodes[report.Node] = incidentcorrelation.ConvertConcreteIncidentArrayIntoIncidents(report.Incidents)
	}

	return s.incidentMerger.ScheduleIncidentsMerge(incidentsByNodes)
}

func (s *ReportsService) getTitleForReport(cluster string, fromDate time.Time, toDate time.Time) string {
	return fmt.Sprintf("Report for %s (%s - %s)",
		cluster,
		fmt.Sprintf("%d.%d", fromDate.Month(), fromDate.Year()),
		fmt.Sprintf("%d.%d", toDate.Month(), toDate.Year()),
	)
}

func (s *ReportsService) getApplicationIncidentFromInsight(
	insight insights.ApplicationInsightsWithMetadata,
	configuration *insights.ApplicationInsightConfiguration,
) *repositories.ApplicationIncident {

	sources := array.Map(func(metadata insights.ApplicationInsightMetadata) repositories.ApplicationIncidentSource {
		return repositories.ApplicationIncidentSource{
			ContainerName: metadata.ContainerName,
			PodName:       metadata.PodName,
			Content:       metadata.Source,
			Timestamp:     metadata.CollectedAtMs,
			Image:         metadata.Image,
		}
	})(insight.Metadata)

	if len(insight.Metadata) == 0 {
		s.logger.Info("Metadata is empty", zap.Any("metadata", insight), zap.Any("conf", configuration))
	}

	return &repositories.ApplicationIncident{
		ApplicationName: insight.Metadata[0].ApplicationName,
		ClusterId:       insight.Metadata[0].ClusterId,
		CustomPrompt:    configuration.CustomPrompt,
		Accuracy:        configuration.Accuracy,
		Title:           insight.Insight.Title,
		Category:        insight.Insight.Category,
		Summary:         insight.Insight.Summary,
		Recommendation:  insight.Insight.Recommendation,
		Urgency:         insight.Insight.Urgency,
		Sources:         sources,
	}
}

// Get maximum of all urgencies from incidents from passed reports
func (s *ReportsService) getReportUrgencyFromApplicationAndNodeReports(
	applicationReports []*repositories.ApplicationReport,
	nodeReports []*repositories.NodeReport,
) insights.Urgency {

	applicationUrgnency := array.Map(func(report *repositories.ApplicationReport) []insights.Urgency {
		return array.Map(func(incident *repositories.ApplicationIncident) insights.Urgency {
			return incident.Urgency
		})(report.Incidents)
	})(applicationReports)

	nodeUrgency := array.Map(func(report *repositories.NodeReport) []insights.Urgency {
		return array.Map(func(incident *repositories.NodeIncident) insights.Urgency {
			return incident.Urgency
		})(report.Incidents)
	})(nodeReports)

	flattenedApplicationUrgency := array.Flatten(applicationUrgnency)
	flattenedNodeUrgency := array.Flatten(nodeUrgency)

	allUrgencies := append(flattenedApplicationUrgency, flattenedNodeUrgency...)

	if len(allUrgencies) == 0 {
		return insights.Urgency_Low
	}

	return insights.MaxUrgency(allUrgencies)
}

func (s *ReportsService) getNodeIncidentFromInsight(insight insights.NodeInsightsWithMetadata, configuration *insights.NodeInsightConfiguration) *repositories.NodeIncident {

	sources := array.Map(func(metadata insights.NodeInsightMetadata) repositories.NodeIncidentSource {
		return repositories.NodeIncidentSource{
			Content:   metadata.Source,
			Timestamp: metadata.CollectedAtMs,
			Filename:  metadata.Filename,
		}
	})(insight.Metadata)

	return &repositories.NodeIncident{
		ClusterId:      insight.Metadata[0].ClusterId,
		NodeName:       insight.Metadata[0].NodeName,
		CustomPrompt:   configuration.CustomPrompt,
		Accuracy:       configuration.Accuracy,
		Title:          insight.Insight.Title,
		Category:       insight.Insight.Category,
		Summary:        insight.Insight.Summary,
		Recommendation: insight.Insight.Recommendation,
		Urgency:        insight.Insight.Urgency,
		Sources:        sources,
	}
}

func (s *ReportsService) GetApplicationReportsFromInsights(
	applicationInsights []insights.ApplicationInsightsWithMetadata,
	applicationConfiguration []*insights.ApplicationInsightConfiguration,
) ([]*repositories.ApplicationReport, error) {

	insightsByApplication := insights.GroupInsightsByApplication(applicationInsights)

	reports := make([]*repositories.ApplicationReport, 0, len(insightsByApplication))
	configByApp := insights.MapApplicationNameToConfiguration(applicationConfiguration)

	for applicationName, insightsForApplication := range insightsByApplication {

		incidentsFromInsights := array.Map(func(insight insights.ApplicationInsightsWithMetadata) *repositories.ApplicationIncident {

			s.logger.Info("config", zap.Any("app", insight.Insight.ApplicationName))
			return s.getApplicationIncidentFromInsight(insight, configByApp[insight.Insight.ApplicationName])
		})

		incidents := incidentsFromInsights(insightsForApplication)
		report := &repositories.ApplicationReport{
			ApplicationName: applicationName,
			Incidents:       incidents,
		}

		config, ok := configByApp[applicationName]
		if ok {
			report.Accuracy = config.Accuracy
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (s *ReportsService) GetNodeReportsFromInsights(
	nodeInsights []insights.NodeInsightsWithMetadata,
	nodesConfiguration []*insights.NodeInsightConfiguration,
) ([]*repositories.NodeReport, error) {
	insightsByNode := insights.GroupInsightsByNode(nodeInsights)

	reports := make([]*repositories.NodeReport, 0, len(insightsByNode))
	configByNode := insights.MapNodeNameToConfiguration(nodesConfiguration)

	for nodeName, insightsForNode := range insightsByNode {

		nodeIncidentsFromInsights := array.Map(func(insight insights.NodeInsightsWithMetadata) *repositories.NodeIncident {
			return s.getNodeIncidentFromInsight(insight, configByNode[insight.Insight.NodeName])
		})

		incidents := nodeIncidentsFromInsights(insightsForNode)

		report := &repositories.NodeReport{
			Node:      nodeName,
			Incidents: incidents,
		}

		config, ok := configByNode[nodeName]
		if ok {
			report.Accuracy = config.Accuracy
			report.CustomPrompt = config.CustomPrompt
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (s *ReportsService) GetApplicationLogsByParams(
	ctx context.Context,
	clusterId string,
	sinceDate time.Time,
	toDate time.Time,
) (sharedrepositories.ApplicationLogsBatchRetriever, error) {
	applicationLogs, err := s.applicationLogsRepository.GetBatchedLogs(ctx,
		clusterId,
		sinceDate,
		toDate)

	if err != nil {
		s.logger.Error("Failed to get application logs", zap.Error(err))
		return nil, err
	}

	return applicationLogs, nil

}

// func (s *ReportsService) GenerateApplicationReports(
// 	ctx context.Context,
// 	applicationLogs []*sharedrepositories.ApplicationLogsDocument,
// 	applicationConfiguration []*insights.ApplicationInsightConfiguration,
// ) ([]*repositories.ApplicationReport, error) {
//
// 	applicationInsights, err := s.applicationInsightsGenerator.OnDemandApplicationInsights(
// 		applicationLogs,
// 		applicationConfiguration,
// 	)
// 	if err != nil {
// 		s.logger.Error("Failed to generate application insights", zap.Error(err))
// 		return nil, err
// 	}
//
// 	return s.GetApplicationReportsFromInsights(
// 		applicationInsights,
// 		applicationConfiguration,
// 	)
// }

func (s *ReportsService) GetNodeLogsByParams(
	ctx context.Context,
	clusterId string,
	fromDate time.Time,
	toDate time.Time,
) (sharedrepositories.NodeLogsBatchRetriever, error) {

	nodeLogs, err := s.nodeLogsRepository.GetBatchedLogs(ctx,
		clusterId,
		fromDate,
		toDate)
	if err != nil {
		s.logger.Error("Failed to get node logs", zap.Error(err))
		return nil, err
	}

	return nodeLogs, nil
}

// func (s *ReportsService) GenerateNodeReports(
//
//	ctx context.Context,
//	nodeLogs []*sharedrepositories.NodeLogsDocument,
//	nodeConfiguration []*insights.NodeInsightConfiguration,
//
// ) ([]*repositories.NodeReport, error) {
//
//		nodeInsights, err := s.nodeInsightsGenerator.OnDemandNodeInsights(
//			nodeLogs,
//			nodeConfiguration,
//		)
//		if err != nil {
//			s.logger.Error("Failed to generate application insights", zap.Error(err))
//			return nil, err
//		}
//
//		return s.GetNodeReportsFromInsights(
//			nodeInsights,
//			nodeConfiguration,
//		)
//	}
func (s *ReportsService) GetSingleReport(ctx context.Context, id string) (*repositories.Report, *repositories.ReportRepositoryError) {
	return s.reportRepository.GetSingleReport(ctx, id)
}

func (s *ReportsService) GetAllReports(ctx context.Context, filter repositories.FilterParams) ([]*repositories.Report, *repositories.ReportRepositoryError) {
	return s.reportRepository.GetAllReports(ctx, filter)
}

func (s *ReportsService) InsertReport(ctx context.Context, report *repositories.Report) (*repositories.Report, *repositories.ReportRepositoryError) {
	return s.reportRepository.InsertReport(ctx, report)
}
