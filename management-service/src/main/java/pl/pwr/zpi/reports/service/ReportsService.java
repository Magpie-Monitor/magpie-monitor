package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
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
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.repository.ApplicationIncidentRepository;
import pl.pwr.zpi.reports.repository.NodeIncidentRepository;
import pl.pwr.zpi.reports.repository.ReportGenerationRequestMetadataRepository;
import pl.pwr.zpi.reports.repository.ReportRepository;

import java.util.List;
import java.util.Optional;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReportsService {

    private final ReportRepository reportRepository;
    private final NodeIncidentRepository nodeIncidentRepository;
    private final ApplicationIncidentRepository applicationIncidentRepository;

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
            return ReportIncidentsDTO.builder()
                    .applicationIncidents(incidents.getApplicationIncidents())
                    .nodeIncidents(incidents.getNodeIncidents())
                    .build();
        });
    }

    public Optional<NodeIncident> getNodeIncidentById(String incidentId) {
        return nodeIncidentRepository.findById(incidentId);
    }


    public Optional<ApplicationIncident> getApplicationIncidentById(String incidentId) {
        return applicationIncidentRepository.findById(incidentId);
    }
}
