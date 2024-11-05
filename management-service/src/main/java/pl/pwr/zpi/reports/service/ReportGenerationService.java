package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.NotificationService;
import pl.pwr.zpi.reports.broker.ReportPublisher;
import pl.pwr.zpi.reports.dto.event.ReportGenerated;
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed;
import pl.pwr.zpi.reports.dto.event.ReportRequested;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.repository.ApplicationIncidentRepository;
import pl.pwr.zpi.reports.repository.NodeIncidentRepository;
import pl.pwr.zpi.reports.repository.ReportGenerationRequestMetadataRepository;
import pl.pwr.zpi.reports.repository.ReportRepository;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReportGenerationService {

    private final ReportPublisher reportPublisher;
    private final NotificationService notificationService;

    private final ReportRepository reportRepository;
    private final NodeIncidentRepository nodeIncidentRepository;
    private final ApplicationIncidentRepository applicationIncidentRepository;

    private final ReportGenerationRequestMetadataRepository reportGenerationRequestMetadataRepository;

    public void retryFailedReportGenerationRequest(String correlationId) {
        reportGenerationRequestMetadataRepository.findByCorrelationId(correlationId).ifPresent(metadata -> {
            createReport(metadata.getCreateReportRequest());
        });
    }
    
    public void createReport(CreateReportRequest reportRequest) {
        ReportRequested reportRequested = ReportRequested.of(reportRequest);
        persistReportGenerationRequestMetadata(reportRequested.correlationId(), reportRequest);
        reportPublisher.publishReportRequestedEvent(reportRequested, this::handleReportGenerationError);
    }

    public void persistReportGenerationRequestMetadata(String correlationId, CreateReportRequest reportRequest) {
        reportGenerationRequestMetadataRepository.save(
                ReportGenerationRequestMetadata.fromCreateReportRequest(correlationId, reportRequest)
        );
    }

    public void handleReportGenerationError(ReportRequestFailed requestFailed) {
        log.error("Report generation request failed: {}", requestFailed);
        reportGenerationRequestMetadataRepository.findByCorrelationId(requestFailed.correlationId())
                .ifPresentOrElse(this::failReportGenerationRequest, () -> {
                    throw new RuntimeException(
                            String.format("Report generation request of correlationId: %s has failed, " +
                                    "but there's no corresponding request metadata.", requestFailed.correlationId()
                            ));
                });
    }

    private void failReportGenerationRequest(ReportGenerationRequestMetadata requestMetadata) {
        log.info("Report generation request failed, correlationId: {}, clusterId: {}", requestMetadata.getCorrelationId(), requestMetadata.getCreateReportRequest().clusterId());

        updateReportGenerationRequestMetadataStatus(requestMetadata, ReportGenerationStatus.ERROR);
        notifyReportGenerationFailed(requestMetadata);
    }

    public void handleReportGenerated(ReportGenerated reportGenerated) {
        log.info("Report generated, correlationId: {}, clusterId: {}", reportGenerated.correlationId(), reportGenerated.report().getClusterId());

        reportGenerationRequestMetadataRepository.findByCorrelationId(reportGenerated.correlationId())
                .ifPresent(requestMetadata -> saveGeneratedReport(requestMetadata, reportGenerated));
    }

    private void saveGeneratedReport(ReportGenerationRequestMetadata requestMetadata, ReportGenerated reportGenerated) {
        persistReport(reportGenerated.report());
        updateReportGenerationRequestMetadataStatus(requestMetadata, ReportGenerationStatus.GENERATED);
        notifyReportGenerated(requestMetadata);
    }

    private void updateReportGenerationRequestMetadataStatus(ReportGenerationRequestMetadata requestMetadata, ReportGenerationStatus generationStatus) {
        requestMetadata.setStatus(generationStatus);
        reportGenerationRequestMetadataRepository.save(requestMetadata);
    }

    private void persistReport(Report report) {
        reportRepository.save(report);
        nodeIncidentRepository.saveAll(report.getNodeIncidents());
        applicationIncidentRepository.saveAll(report.getApplicationIncidents());
    }

    // TODO - stub implementation
    public void notifyReportGenerated(ReportGenerationRequestMetadata requestMetadata) {
        notificationService.notifySlack(requestMetadata.getSlackReceiverIds());
        notificationService.notifyDiscord(requestMetadata.getDiscordReceiverIds());
        notificationService.notifyEmail(requestMetadata.getMailReceiverIds());
    }

    // TODO - stub implementation
    public void notifyReportGenerationFailed(ReportGenerationRequestMetadata requestMetadata) {
        notificationService.notifySlack(requestMetadata.getSlackReceiverIds());
        notificationService.notifyDiscord(requestMetadata.getDiscordReceiverIds());
        notificationService.notifyEmail(requestMetadata.getMailReceiverIds());
    }
}
