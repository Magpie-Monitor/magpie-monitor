package pl.pwr.zpi.reports

import pl.pwr.zpi.reports.dto.event.ReportGenerated
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed
import pl.pwr.zpi.reports.dto.event.ReportRequested
import pl.pwr.zpi.reports.dto.request.CreateReportRequest
import pl.pwr.zpi.reports.entity.report.Report
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata
import pl.pwr.zpi.reports.enums.Accuracy
import pl.pwr.zpi.reports.enums.ReportGenerationStatus
import pl.pwr.zpi.reports.enums.ReportType
import pl.pwr.zpi.reports.service.ReportGenerationService
import pl.pwr.zpi.reports.repository.*
import pl.pwr.zpi.notifications.ReportNotificationService
import pl.pwr.zpi.reports.broker.ReportPublisher
import spock.lang.Specification

class ReportGenerationServiceTest extends Specification {

    def reportPublisher = Mock(ReportPublisher)
    def reportNotificationService = Mock(ReportNotificationService)
    def reportRepository = Mock(ReportRepository)
    def nodeIncidentRepository = Mock(NodeIncidentRepository)
    def nodeIncidentSourcesRepository = Mock(NodeIncidentSourcesRepository)
    def applicationIncidentRepository = Mock(ApplicationIncidentRepository)
    def applicationIncidentSourcesRepository = Mock(ApplicationIncidentSourcesRepository)
    def reportGenerationRequestMetadataRepository = Mock(ReportGenerationRequestMetadataRepository)

    def reportGenerationService = new ReportGenerationService(
            reportPublisher,
            reportNotificationService,
            reportRepository,
            nodeIncidentRepository,
            nodeIncidentSourcesRepository,
            applicationIncidentRepository,
            applicationIncidentSourcesRepository,
            reportGenerationRequestMetadataRepository
    )

    def "should create report and publish report requested event"() {
        given:
        def reportRequest = createCreateReportRequest("cluster123", 0L, 86400000L)
        def reportType = ReportType.SCHEDULED
        def reportRequested = ReportRequested.of(reportRequest)

        when:
        reportGenerationService.createReport(reportRequest, reportType)

        then:
        1 * reportGenerationRequestMetadataRepository.save(_ as ReportGenerationRequestMetadata)
        1 * reportPublisher.publishReportRequestedEvent({ ReportRequested reportRequestedEvent ->
            reportRequestedEvent.reportRequest.clusterId == reportRequested.reportRequest.clusterId &&
                    reportRequestedEvent.reportRequest.sinceMs == reportRequested.reportRequest.sinceMs &&
                    reportRequestedEvent.reportRequest.toMs == reportRequested.reportRequest.toMs &&
                    reportRequestedEvent.reportRequest.applicationConfiguration == reportRequested.reportRequest.applicationConfiguration &&
                    reportRequestedEvent.reportRequest.nodeConfiguration == reportRequested.reportRequest.nodeConfiguration
        }, _)
    }

    def "should handle report generation failure and notify"() {
        given:
        def correlationId = "correlation123"
        def requestFailed = ReportRequestFailed.builder()
                .correlationId(correlationId)
                .errorType("Failed to generate report")
                .errorMessage("Error occurred")
                .timestampMs(System.currentTimeMillis())
                .build()
        def createReportRequest = CreateReportRequest.builder()
                .clusterId("cluster123")
                .accuracy(Accuracy.HIGH)
                .sinceMs(0L)
                .toMs(86400000L)
                .slackReceiverIds([])
                .emailReceiverIds([])
                .discordReceiverIds([])
                .applicationConfigurations([])
                .nodeConfigurations([])
                .build()
        def reportMetadata = ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.ERROR)
                .error(requestFailed)
                .reportType(ReportType.SCHEDULED)
                .createReportRequest(createReportRequest)
                .build()

        when:
        reportGenerationService.handleReportGenerationError(requestFailed)

        then:
        1 * reportGenerationRequestMetadataRepository.findByCorrelationId(correlationId) >> Optional.of(reportMetadata)
        1 * reportNotificationService.notifySlackOnReportGenerationFailed(_, _)
        1 * reportNotificationService.notifyDiscordOnReportGenerationFailed(_, _)
        1 * reportNotificationService.notifyEmailOnReportGenerationFailed(_, _)
    }

    def "should handle report generated and save report"() {
        given:
        def correlationId = "correlation123"
        def report = Report.builder()
                .nodeReports([])
                .applicationReports([])
                .build()
        def reportGenerated = new ReportGenerated(correlationId, report, System.currentTimeMillis())
        def createReportRequest = CreateReportRequest.builder()
                .clusterId("cluster123")
                .accuracy(Accuracy.HIGH)
                .sinceMs(0L)
                .toMs(86400000L)
                .slackReceiverIds([])
                .emailReceiverIds([])
                .discordReceiverIds([])
                .applicationConfigurations([])
                .nodeConfigurations([])
                .build()
        def reportMetadata = ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATED)
                .reportType(ReportType.SCHEDULED)
                .createReportRequest(createReportRequest)
                .build()

        when:
        reportGenerationService.handleReportGenerated(reportGenerated)

        then:
        1 * reportGenerationRequestMetadataRepository.findByCorrelationId(correlationId) >> Optional.of(reportMetadata)
        1 * reportRepository.save(_ as Report)
        1 * reportNotificationService.notifySlackOnReportCreated(_, _)
        1 * reportNotificationService.notifyDiscordOnReportCreated(_, _)
        1 * reportNotificationService.notifyEmailOnReportCreated(_, _)
    }

    def "should retry failed report generation request"() {
        given:
        def correlationId = "correlation123"
        def reportRequest = createCreateReportRequest("cluster123", 0L, 86400000L)
        def reportMetadata = ReportGenerationRequestMetadata.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATED)
                .reportType(ReportType.SCHEDULED)
                .createReportRequest(reportRequest)
                .build()

        reportGenerationRequestMetadataRepository.findByCorrelationId(correlationId) >> Optional.of(reportMetadata)

        when:
        reportGenerationService.retryFailedReportGenerationRequest(correlationId)

        then:
        1 * reportPublisher.publishReportRequestedEvent(_, _)
    }

    def "should save report generation pl.pwr.zpi.metadata"() {
        given:
        def correlationId = "correlation123"
        def reportRequest = createCreateReportRequest("cluster123", 0L, 86400000L)
        def reportType = ReportType.SCHEDULED

        when:
        reportGenerationService.persistReportGenerationRequestMetadata(correlationId, reportRequest, reportType)

        then:
        1 * reportGenerationRequestMetadataRepository.save(_ as ReportGenerationRequestMetadata)
    }

    def "should throw exception if no pl.pwr.zpi.metadata found on report generation failure"() {
        given:
        def correlationId = "correlation123"
        def requestFailed = ReportRequestFailed.builder()
                .correlationId(correlationId)
                .errorType("Failed to generate report")
                .errorMessage("Error occurred")
                .timestampMs(System.currentTimeMillis())
                .build()

        reportGenerationRequestMetadataRepository.findByCorrelationId(correlationId) >> Optional.empty()

        when:
        reportGenerationService.handleReportGenerationError(requestFailed)

        then:
        thrown(RuntimeException)
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
}
