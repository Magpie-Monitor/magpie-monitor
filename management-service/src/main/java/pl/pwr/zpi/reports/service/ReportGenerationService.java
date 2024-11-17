package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.notifications.ReportNotificationService;
import pl.pwr.zpi.reports.broker.ReportPublisher;
import pl.pwr.zpi.reports.dto.event.ReportGenerated;
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed;
import pl.pwr.zpi.reports.dto.event.ReportRequested;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.dto.request.CreateReportScheduleRequest;
import pl.pwr.zpi.reports.dto.scheduler.ReportSchedule;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.repository.*;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReportGenerationService {

    private final ReportPublisher reportPublisher;
    private final ReportNotificationService reportNotificationService;
    private final ReportRepository reportRepository;
    private final NodeIncidentRepository nodeIncidentRepository;
    private final NodeIncidentSourcesRepository nodeIncidentSourcesRepository;
    private final ApplicationIncidentRepository applicationIncidentRepository;
    private final ApplicationIncidentSourcesRepository applicationIncidentSourcesRepository;
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
                .ifPresentOrElse(
                        metadata -> failReportGenerationRequest(metadata, requestFailed),
                        () -> {
                            throw new RuntimeException(
                                    String.format("Report generation request of correlationId: %s has failed, " +
                                            "but there's no corresponding request metadata.", requestFailed.correlationId()
                                    ));
                        });
    }

    private void failReportGenerationRequest(ReportGenerationRequestMetadata requestMetadata,
                                             ReportRequestFailed requestFailed) {
        log.info("Report generation request failed, correlationId: {}, clusterId: {}", requestMetadata.getCorrelationId(), requestMetadata.getCreateReportRequest().clusterId());

        markReportGenerationRequestAsFailed(requestMetadata, requestFailed);
        notifyReportGenerationFailed(requestMetadata, requestMetadata.getClusterId());
    }

    private void markReportGenerationRequestAsFailed(
            ReportGenerationRequestMetadata requestMetadata, ReportRequestFailed requestFailed) {
        requestMetadata.setError(requestFailed);
        updateReportGenerationRequestMetadataStatus(requestMetadata, ReportGenerationStatus.ERROR);
    }

    public void handleReportGenerated(ReportGenerated reportGenerated) {
        log.info("Report generated, correlationId: {}, clusterId: {}", reportGenerated.correlationId(), reportGenerated.report().getClusterId());

        reportGenerationRequestMetadataRepository.findByCorrelationId(reportGenerated.correlationId())
                .ifPresent(requestMetadata -> saveGeneratedReport(requestMetadata, reportGenerated));
    }

    private void saveGeneratedReport(ReportGenerationRequestMetadata requestMetadata, ReportGenerated reportGenerated) {
        persistReport(reportGenerated.report());
        updateReportGenerationRequestMetadataStatus(requestMetadata, ReportGenerationStatus.GENERATED);
        notifyReportGenerated(requestMetadata, reportGenerated.getReportId());
    }

    private void updateReportGenerationRequestMetadataStatus(ReportGenerationRequestMetadata requestMetadata, ReportGenerationStatus generationStatus) {
        requestMetadata.setStatus(generationStatus);
        reportGenerationRequestMetadataRepository.save(requestMetadata);
    }

    private void persistReport(Report report) {
        reportRepository.save(report);

        extendIncidentsWithReportId(report);

        nodeIncidentRepository.saveAll(report.getNodeIncidents());
        applicationIncidentRepository.saveAll(report.getApplicationIncidents());

        applicationIncidentSourcesRepository.saveAll(report.getApplicationIncidentSources());
        nodeIncidentSourcesRepository.saveAll(report.getNodeIncidentSources());
    }

    private void extendIncidentsWithReportId(Report report) {
        report.getNodeIncidents().forEach(
                nodeIncident -> {
                    nodeIncident.setReportId(report.getId());
                    nodeIncident.extendSourcesWithIncidentId();
                }
        );
        report.getApplicationIncidents().forEach(
                applicationIncident -> {
                    applicationIncident.setReportId(report.getId());
                    applicationIncident.extendSourcesWithIncidentId();
                }
        );
    }

    public void notifyReportGenerated(ReportGenerationRequestMetadata requestMetadata, String reportId) {
        reportNotificationService.notifySlackOnReportCreated(requestMetadata.getSlackReceiverIds(), reportId);
        reportNotificationService.notifyDiscordOnReportCreated(requestMetadata.getDiscordReceiverIds(), reportId);
        reportNotificationService.notifyEmailOnReportCreated(requestMetadata.getMailReceiverIds(), reportId);
    }

    public void notifyReportGenerationFailed(ReportGenerationRequestMetadata requestMetadata, String clusterId) {
        reportNotificationService.notifySlackOnReportGenerationFailed(requestMetadata.getSlackReceiverIds(), clusterId);
        reportNotificationService.notifyDiscordOnReportGenerationFailed(requestMetadata.getDiscordReceiverIds(), clusterId);
        reportNotificationService.notifyEmailOnReportGenerationFailed(requestMetadata.getMailReceiverIds(), clusterId);
    }
}
