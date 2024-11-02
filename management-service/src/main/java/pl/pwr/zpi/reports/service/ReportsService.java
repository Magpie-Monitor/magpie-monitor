package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.notifications.NotificationService;
import pl.pwr.zpi.reports.broker.ReportPublisher;
import pl.pwr.zpi.reports.dto.event.ReportGenerated;
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed;
import pl.pwr.zpi.reports.dto.event.ReportRequested;
import pl.pwr.zpi.reports.dto.report.ReportDetailedSummaryDTO;
import pl.pwr.zpi.reports.dto.report.ReportIncidentsDTO;
import pl.pwr.zpi.reports.dto.report.ReportSummaryDTO;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.repository.ApplicationIncidentRepository;
import pl.pwr.zpi.reports.repository.NodeIncidentRepository;
import pl.pwr.zpi.reports.repository.ReportGenerationRequestMetadataRepository;
import pl.pwr.zpi.reports.repository.ReportRepository;
import pl.pwr.zpi.reports.repository.projection.ReportIncidentsProjection;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ReportsService {

    private final ReportPublisher reportPublisher;
    private final NotificationService notificationService;

    private final ReportRepository reportRepository;
    private final ReportGenerationRequestMetadataRepository reportGenerationRequestMetadataRepository;
    private final NodeIncidentRepository nodeIncidentRepository;
    private final ApplicationIncidentRepository applicationIncidentRepository;

    public void createReport(CreateReportRequest reportRequest) {
        ReportRequested reportRequested = ReportRequested.of(reportRequest);
        reportPublisher.publishReportRequestedEvent(reportRequested);
        persistReportGenerationRequestMetadata(reportRequested.correlationId(), reportRequest);
    }

    public void persistReportGenerationRequestMetadata(String correlationId, CreateReportRequest reportRequest) {
        reportGenerationRequestMetadataRepository.save(
                ReportGenerationRequestMetadata.fromCreateReportRequest(correlationId, reportRequest)
        );
    }

    public void handleReportGenerationError(ReportRequestFailed requestFailed) {
        reportGenerationRequestMetadataRepository.findByCorrelationId(requestFailed.correlationId())
                .ifPresent(requestMetadata -> {
                    requestMetadata.markAsFailed();
                    reportGenerationRequestMetadataRepository.save(requestMetadata);
                    notifyReportGenerationFailed(requestMetadata);
                });
    }

    public void handleReportGenerated(ReportGenerated reportGenerated) {
        reportGenerationRequestMetadataRepository.findByCorrelationId(reportGenerated.correlationId())
                .ifPresent(requestMetadata -> {
                    requestMetadata.markAsGenerated();
                    Report report = reportGenerated.report();
                    persistReport(report);
                    persistReportIncidents(report);
                    notifyReportGenerated(requestMetadata);
                });
    }

    private void persistReport(Report report) {
        reportRepository.save(report);
    }

    private void persistReportIncidents(Report report) {
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

    public List<ReportSummaryDTO> getReportSummaries() {
        return reportRepository.findAllProjectedBy().stream()
                .map(ReportSummaryDTO::ofReportSummaryProjection)
                .toList();
    }

    public Optional<ReportDetailedSummaryDTO> getReportDetailedSummaryById(String reportId) {
        return reportRepository.findProjectedBy(reportId)
                .map(ReportDetailedSummaryDTO::fromReportDetailedSummaryProjection);
    }

    public Optional<ReportIncidentsDTO> getReportIncidents(String reportId) {
        return reportRepository.findProjectedById(reportId).map(incidents -> {
            List<ApplicationIncident> applicationIncidents = incidents.getApplicationReports().stream()
                    .map(ReportIncidentsProjection.ApplicationReportProjection::getIncidents)
                    .flatMap(List::stream)
                    .toList();

            List<NodeIncident> nodeIncidents = incidents.getNodeReports().stream()
                    .map(ReportIncidentsProjection.NodeReportProjection::getIncidents)
                    .flatMap(List::stream)
                    .toList();

            return new ReportIncidentsDTO(applicationIncidents, nodeIncidents);
        });
    }

    public Optional<NodeIncident> getNodeIncidentById(String incidentId) {
        return nodeIncidentRepository.findById(incidentId);
    }


    public Optional<ApplicationIncident> getApplicationIncidentById(String incidentId) {
        return applicationIncidentRepository.findById(incidentId);
    }
}
