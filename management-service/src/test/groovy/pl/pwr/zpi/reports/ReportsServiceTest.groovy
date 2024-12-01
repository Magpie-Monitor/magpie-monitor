package pl.pwr.zpi.reports

import pl.pwr.zpi.reports.dto.report.*
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentDTO
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentDTO
import pl.pwr.zpi.reports.dto.request.CreateReportRequest
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource
import pl.pwr.zpi.reports.entity.report.node.NodeIncident
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata
import pl.pwr.zpi.reports.enums.Accuracy
import pl.pwr.zpi.reports.enums.ReportGenerationStatus
import pl.pwr.zpi.reports.enums.ReportType
import pl.pwr.zpi.reports.enums.Urgency
import pl.pwr.zpi.reports.repository.*
import pl.pwr.zpi.reports.repository.projection.ReportDetailedSummaryProjection
import org.springframework.data.domain.Pageable
import pl.pwr.zpi.reports.repository.projection.ReportSummaryProjection
import pl.pwr.zpi.reports.service.ReportsService
import spock.lang.Specification


class ReportsServiceTest extends Specification {

    ReportRepository reportRepository
    NodeIncidentRepository nodeIncidentRepository
    ApplicationIncidentRepository applicationIncidentRepository
    ApplicationIncidentSourcesRepository applicationIncidentSourcesRepository
    NodeIncidentSourcesRepository nonNodeIncidentSourcesRepository
    ReportGenerationRequestMetadataRepository reportGenerationRequestMetadataRepository

    ReportsService reportsService

    def setup() {
        reportRepository = Mock()
        nodeIncidentRepository = Mock()
        applicationIncidentRepository = Mock()
        applicationIncidentSourcesRepository = Mock()
        nonNodeIncidentSourcesRepository = Mock()
        reportGenerationRequestMetadataRepository = Mock()
        reportsService = new ReportsService(reportRepository, nodeIncidentRepository,
                applicationIncidentRepository, applicationIncidentSourcesRepository,
                nonNodeIncidentSourcesRepository, reportGenerationRequestMetadataRepository)
    }

    def "should get failed report generation requests"() {
        given:
        def failedRequests = [Mock(ReportGenerationRequestMetadata), Mock(ReportGenerationRequestMetadata)]
        reportGenerationRequestMetadataRepository.findByStatus(ReportGenerationStatus.ERROR) >> failedRequests

        when:
        def result = reportsService.getFailedReportGenerationRequests()

        then:
        result == failedRequests
    }

    def "should get generation reports"() {
        given:
        def generatingReports = List.of(
                createReportGenerationRequestMetadata("132321", ReportGenerationStatus.GENERATING, "cluster1", 1000L, 2000L),
                createReportGenerationRequestMetadata("132322", ReportGenerationStatus.GENERATING, "cluster2", 1000L, 2000L))
        reportGenerationRequestMetadataRepository.findByStatus(ReportGenerationStatus.GENERATING) >> generatingReports

        when:
        def result = reportsService.getGenerationReports()

        then:
        result.size() == generatingReports.size()
    }

    def "should get report summaries"() {
        given:
        def reportSummary = Mock(ReportSummaryProjection) {
            getId() >> "123"
            getClusterId() >> "cluster1"
            getTitle() >> "Scheduled Report"
            getUrgency() >> Urgency.HIGH
            getRequestedAtMs() >> System.currentTimeMillis()
            getSinceMs() >> System.currentTimeMillis() - 10000L
            getToMs() >> System.currentTimeMillis()
        }

        def reportSummaries = [reportSummary]
        reportRepository.findAllByReportType(ReportType.SCHEDULED) >> reportSummaries

        def expectedReportSummaries = reportSummaries.collect { report ->
            ReportSummaryDTO.builder()
                    .id(report.id)
                    .clusterId(report.clusterId)
                    .title(report.title)
                    .urgency(report.urgency)
                    .requestedAtMs(report.requestedAtMs)
                    .sinceMs(report.sinceMs)
                    .toMs(report.toMs)
                    .build()
        }

        when:
        def result = reportsService.getReportSummaries("SCHEDULED")

        then:
        result == expectedReportSummaries
    }

    def "should get report detailed summary by id"() {
        given:
        def reportId = "report123"

        def reportDetailedSummaryProjection = Mock(ReportDetailedSummaryProjection) {
            getId() >> reportId
            getClusterId() >> "cluster1"
            getTitle() >> "Detailed Report"
            getUrgency() >> Urgency.HIGH
            getRequestedAtMs() >> 1732461533844L
            getSinceMs() >> 1732461523846L
            getToMs() >> 1732461533877L
            getTotalApplicationEntries() >> 10
            getTotalNodeEntries() >> 5
            getAnalyzedApplications() >> 8
            getAnalyzedNodes() >> 4
        }

        reportRepository.findProjectedDetailedById(reportId) >> Optional.of(reportDetailedSummaryProjection)

        when:
        def result = reportsService.getReportDetailedSummaryById(reportId)

        then:
        result.isPresent()
        result.get().id == reportId
        result.get().clusterId == "cluster1"
        result.get().title == "Detailed Report"
        result.get().urgency == Urgency.HIGH
        result.get().requestedAtMs == 1732461533844L
        result.get().sinceMs == 1732461523846L
        result.get().toMs == 1732461533877L
        result.get().totalApplicationEntries == 10
        result.get().totalNodeEntries == 5
        result.get().analyzedApplications == 8
        result.get().analyzedNodes == 4
    }


    def "should return empty optional when report detailed summary not found"() {
        given:
        def reportId = "report123"
        reportRepository.findProjectedDetailedById(reportId) >> Optional.empty()

        when:
        def result = reportsService.getReportDetailedSummaryById(reportId)

        then:
        result == Optional.empty()
    }

    def "should get node incidents by report id"() {
        given:
        def reportId = "report123"
        def pageable = Pageable.unpaged()

        def nodeIncidents = [createNodeIncident("incident123", reportId)]
        def nodeIncidentDTOs = nodeIncidents.collect { NodeIncidentDTO.fromNodeIncident(it) }

        nodeIncidentRepository.findByReportId(reportId, pageable) >> nodeIncidents
        nodeIncidentRepository.countByReportId(reportId) >> nodeIncidents.size()

        when:
        def result = reportsService.getReportNodeIncidents(reportId, pageable)

        then:
        result.data == nodeIncidentDTOs
        result.totalEntries == nodeIncidents.size()
    }

    def "should get application incidents by report id"() {
        given:
        def reportId = "report123"
        def pageable = Pageable.unpaged()

        def applicationIncidents = List.of(createApplicationIncident("incident123", reportId))
        def applicationIncidentDTOs = applicationIncidents.collect { ApplicationIncidentDTO.fromApplicationIncident(it) }

        applicationIncidentRepository.findByReportId(reportId, pageable) >> applicationIncidents
        applicationIncidentRepository.countByReportId(reportId) >> applicationIncidents.size()

        when:
        def result = reportsService.getReportApplicationIncidents(reportId, pageable)

        then:
        result.data == applicationIncidentDTOs
        result.totalEntries == applicationIncidents.size()
    }

    def "should get application incident sources by incident id"() {
        given:
        def incidentId = "incident123"
        def pageable = Mock(Pageable)
        def sources = [Mock(ApplicationIncidentSource)]
        applicationIncidentSourcesRepository.findByIncidentId(incidentId, pageable) >> sources
        applicationIncidentSourcesRepository.countByIncidentId(incidentId) >> sources.size()

        when:
        def result = reportsService.getApplicationIncidentSourcesByIncidentId(incidentId, pageable)

        then:
        result.data == sources
        result.totalEntries == sources.size()
    }

    def "should get node incident sources by incident id"() {
        given:
        def incidentId = "incident123"
        def pageable = Mock(Pageable)
        def sources = [Mock(NodeIncidentSource)]
        nonNodeIncidentSourcesRepository.findByIncidentId(incidentId, pageable) >> sources
        nonNodeIncidentSourcesRepository.countByIncidentId(incidentId) >> sources.size()

        when:
        def result = reportsService.getNodeIncidentSourcesByIncidentId(incidentId, pageable)

        then:
        result.data == sources
        result.totalEntries == sources.size()
    }

    def "should get application incident by id"() {
        given:
        def incidentId = "incident123"
        def applicationIncident = createApplicationIncident("test123", incidentId)
        def applicationIncidentDTO = ApplicationIncidentDTO.fromApplicationIncident(applicationIncident)

        applicationIncidentRepository.findById(incidentId) >> Optional.of(applicationIncident)

        when:
        def result = reportsService.getApplicationIncidentById(incidentId)

        then:
        result == Optional.of(applicationIncidentDTO)
    }

    def "should get node incident by id"() {
        given:
        def incidentId = "incident123"
        def nodeIncident = createNodeIncident("test123", incidentId)
        def nodeIncidentDTO = NodeIncidentDTO.fromNodeIncident(nodeIncident)

        nodeIncidentRepository.findById(incidentId) >> Optional.of(nodeIncident)

        when:
        def result = reportsService.getNodeIncidentById(incidentId)

        then:
        result == Optional.of(nodeIncidentDTO)
    }

    private ReportGenerationRequestMetadata createReportGenerationRequestMetadata(String correlationId, ReportGenerationStatus status, String clusterId, long sinceMs, long toMs) {
        return ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(status)
                .createReportRequest(createCreateReportRequest(clusterId, sinceMs, toMs))
                .reportType(ReportType.SCHEDULED)
                .build()
    }

    private CreateReportRequest createCreateReportRequest(String clusterId, long sinceMs, long toMs) {
        return CreateReportRequest.builder()
                .clusterId(clusterId)
                .accuracy(Accuracy.HIGH)
                .sinceMs(sinceMs)
                .toMs(toMs)
                .slackReceiverIds([])
                .emailReceiverIds([])
                .discordReceiverIds([])
                .applicationConfigurations([])
                .nodeConfigurations([])
                .build()
    }

    private ApplicationIncident createApplicationIncident(String id, String reportId) {
        return ApplicationIncident.builder()
                .id(id)
                .reportId(reportId)
                .title("Test Incident")
                .accuracy(Accuracy.HIGH)
                .customPrompt("Custom prompt")
                .clusterId("cluster123")
                .applicationName("TestApp")
                .category("Category1")
                .summary("Incident Summary")
                .recommendation("Recommendation for incident")
                .urgency(Urgency.HIGH)
                .sources([])
                .build()
    }

    private NodeIncident createNodeIncident(String incidentId, String reportId) {
        return NodeIncident.builder()
                .id(incidentId)
                .reportId(reportId)
                .title("Test Node Incident")
                .clusterId("cluster123")
                .nodeName("Node1")
                .category("Node Category")
                .summary("Node Incident Summary")
                .recommendation("Recommendation for node incident")
                .urgency(Urgency.HIGH)
                .sources([])
                .build()
    }
}
