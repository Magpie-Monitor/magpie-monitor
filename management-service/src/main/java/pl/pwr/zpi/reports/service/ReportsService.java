package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.dto.report.ReportDetailedSummaryDTO;
import pl.pwr.zpi.reports.dto.report.ReportIncidentsDTO;
import pl.pwr.zpi.reports.dto.report.ReportSummaryDTO;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentDTO;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.repository.ApplicationIncidentRepository;
import pl.pwr.zpi.reports.repository.NodeIncidentRepository;
import pl.pwr.zpi.reports.repository.ReportGenerationRequestMetadataRepository;
import pl.pwr.zpi.reports.repository.ReportRepository;
import pl.pwr.zpi.reports.repository.projection.ReportIncidentsProjection;

import java.util.List;
import java.util.Optional;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReportsService {

    private final ReportRepository reportRepository;
    private final NodeIncidentRepository nodeIncidentRepository;
    private final ApplicationIncidentRepository applicationIncidentRepository;
    private final ReportGenerationRequestMetadataRepository reportGenerationRequestMetadataRepository;

    public List<ReportGenerationRequestMetadata> getFailedReportGenerationRequests() {
        return reportGenerationRequestMetadataRepository.findByStatus(ReportGenerationStatus.ERROR);
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
        return reportRepository.findProjectedById(reportId).map(incidentProjection -> {
            return ReportIncidentsDTO.builder()
                    .applicationIncidents(extractApplicationIncidents(incidentProjection))
                    .nodeIncidents(extractNodeIncidents(incidentProjection))
                    .build();
        });
    }

    private List<ApplicationIncident> extractApplicationIncidents(ReportIncidentsProjection reportIncidentsProjection) {
        return reportIncidentsProjection.getApplicationReports().stream()
                .map(ReportIncidentsProjection.ApplicationReportProjection::getIncidents)
                .flatMap(List::stream)
                .toList();
    }

    private List<NodeIncident> extractNodeIncidents(ReportIncidentsProjection reportIncidentsProjection) {
        return reportIncidentsProjection.getNodeReports().stream()
                .map(ReportIncidentsProjection.NodeReportProjection::getIncidents)
                .flatMap(List::stream)
                .toList();
    }

    public Optional<ApplicationIncidentDTO> getApplicationIncidentById(String incidentId) {
        return applicationIncidentRepository.findById(incidentId).map(ApplicationIncidentDTO::fromApplicationIncident);
    }

    public Optional<NodeIncidentDTO> getNodeIncidentById(String incidentId) {
        return nodeIncidentRepository.findById(incidentId).map(NodeIncidentDTO::fromNodeIncident);
    }
}
