package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.dto.report.*;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentDTO;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.enums.ReportType;
import pl.pwr.zpi.reports.repository.*;
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
    private final ApplicationIncidentSourcesRepository applicationIncidentSourcesRepository;
    private final NodeIncidentSourcesRepository nonNodeIncidentSourcesRepository;
    private final ReportGenerationRequestMetadataRepository reportGenerationRequestMetadataRepository;

    public List<ReportGenerationRequestMetadata> getFailedReportGenerationRequests() {
        return reportGenerationRequestMetadataRepository.findByStatus(ReportGenerationStatus.ERROR);
    }

    public List<ReportGeneratingDTO> getGenerationReports() {
        return reportGenerationRequestMetadataRepository
                .findByStatus(ReportGenerationStatus.GENERATING)
                .stream()
                .map(ReportGeneratingDTO::ofReportGenerationRequestMetadata)
                .toList();
    }

    public List<ReportSummaryDTO> getReportSummaries(ReportType reportType) {
        return reportRepository.findAllByReportType(reportType).stream()
                .map(ReportSummaryDTO::ofReportSummaryProjection)
                .toList();
    }

    public Optional<ReportDetailedSummaryDTO> getReportDetailedSummaryById(String reportId) {
        return reportRepository.findProjectedDetailedById(reportId)
                .map(ReportDetailedSummaryDTO::fromReportDetailedSummaryProjection);
    }

    public Optional<ReportIncidentsDTO> getReportIncidents(String reportId) {
        return reportRepository.findProjectedIncidentsById(reportId).map(incidentProjection -> {
            return ReportIncidentsDTO.builder()
                    .applicationIncidents(extractApplicationIncidents(incidentProjection))
                    .nodeIncidents(extractNodeIncidents(incidentProjection))
                    .build();
        });
    }

    public ReportPaginatedIncidentsDTO<ApplicationIncidentDTO> getReportApplicationIncidents(
            String reportId, Pageable pageable) {

        return ReportPaginatedIncidentsDTO.<ApplicationIncidentDTO>builder()
                .data(
                        applicationIncidentRepository.findByReportId(reportId, pageable).stream()
                                .map(ApplicationIncidentDTO::fromApplicationIncident)
                                .toList()
                )
                .totalEntries(applicationIncidentRepository.countByReportId(reportId))
                .build();
    }

    public ReportPaginatedIncidentsDTO<NodeIncidentDTO> getReportNodeIncidents(
            String reportId, Pageable pageable) {

        return ReportPaginatedIncidentsDTO.<NodeIncidentDTO>builder()
                .data(
                        nodeIncidentRepository.findByReportId(reportId, pageable).stream()
                                .map(NodeIncidentDTO::fromNodeIncident)
                                .toList()
                )
                .totalEntries(nodeIncidentRepository.countByReportId(reportId))
                .build();
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

    public ReportPaginatedIncidentsDTO<ApplicationIncidentSource> getApplicationIncidentSourcesByIncidentId(
            String incidentId, Pageable pageable) {
        return ReportPaginatedIncidentsDTO.<ApplicationIncidentSource>builder()
                .data(applicationIncidentSourcesRepository.findByIncidentId(incidentId, pageable))
                .totalEntries(applicationIncidentSourcesRepository.countByIncidentId(incidentId))
                .build();
    }

    public ReportPaginatedIncidentsDTO<NodeIncidentSource> getNodeIncidentSourcesByIncidentId(
            String incidentId, Pageable pageable) {
        return ReportPaginatedIncidentsDTO.<NodeIncidentSource>builder()
                .data(nonNodeIncidentSourcesRepository.findByIncidentId(incidentId, pageable))
                .totalEntries(nonNodeIncidentSourcesRepository.countByIncidentId(incidentId))
                .build();

    }

    public Optional<ApplicationIncidentDTO> getApplicationIncidentById(String incidentId) {
        return applicationIncidentRepository.findById(incidentId).map(ApplicationIncidentDTO::fromApplicationIncident);
    }

    public Optional<NodeIncidentDTO> getNodeIncidentById(String incidentId) {
        return nodeIncidentRepository.findById(incidentId).map(NodeIncidentDTO::fromNodeIncident);
    }
}
