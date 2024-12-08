package pl.pwr.zpi.reports.api;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.reports.dto.report.*;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentDTO;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncidentSimplifiedDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeIncidentSimplifiedDTO;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.dto.request.CreateReportScheduleRequest;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource;
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportType;
import pl.pwr.zpi.reports.service.ReportGenerationService;
import pl.pwr.zpi.reports.service.ReportScheduleService;
import pl.pwr.zpi.reports.service.ReportsService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/reports")
public class ReportsController {

    private final ReportsService reportsService;
    private final ReportGenerationService reportGenerationService;
    private final ReportScheduleService reportScheduleService;

    @PostMapping
    public ResponseEntity<Void> createReport(@RequestBody CreateReportRequest reportRequest) {
        reportGenerationService.createReport(reportRequest, ReportType.ON_DEMAND);
        return ResponseEntity.ok().build();
    }

    @PostMapping("/schedule")
    public ResponseEntity<Void> scheduleReport(@RequestBody @Valid CreateReportScheduleRequest reportRequest) {
        reportScheduleService.scheduleReport(reportRequest);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/requests/failed")
    public ResponseEntity<List<ReportGenerationRequestMetadata>> getFailedReportGenerationRequests() {
        return ResponseEntity.ok(reportsService.getFailedReportGenerationRequests());
    }

    @PostMapping("/requests/{id}/retry")
    public ResponseEntity<Void> retryFailedReportGenerationRequest(@PathVariable String id) {
        reportGenerationService.retryFailedReportGenerationRequest(id);
        return ResponseEntity.ok().build();
    }

    @GetMapping
    public ResponseEntity<List<ReportSummaryDTO>> getReportSummaries(@RequestParam String reportType) {
        return ResponseEntity.ok().body(reportsService.getReportSummaries(reportType));
    }

    @GetMapping("/await-generation")
    public ResponseEntity<List<ReportGeneratingDTO>> getAwaitingGenerationReports() {
        return ResponseEntity.ok(reportsService.getAwaitingGenerationReports());
    }

    @GetMapping("/{id}")
    public ResponseEntity<ReportDetailedSummaryDTO> getReportById(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getReportDetailedSummaryById(id));
    }

    @GetMapping("/latest")
    public ResponseEntity<ReportDetailedSummaryDTO> getNewestReport() {
        return ResponseEntity.of(reportsService.getLatestReportDetailedSummary());
    }

    @GetMapping("/{id}/incidents")
    public ResponseEntity<ReportIncidentsDTO> getReportIncidents(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getReportIncidents(id));
    }

    @GetMapping("/{id}/application-incidents")
    public ResponseEntity<ReportPaginatedIncidentsDTO<ApplicationIncidentDTO>> getApplicationIncidentsForReport(
            @PathVariable String id, Pageable pageable
    ) {
        return ResponseEntity.ok(reportsService.getReportApplicationIncidents(id, pageable));
    }

    @GetMapping("/{id}/node-incidents")
    public ResponseEntity<ReportPaginatedIncidentsDTO<NodeIncidentDTO>> getNodeIncidentsForReport(
            @PathVariable String id, Pageable pageable
    ) {
        return ResponseEntity.ok(reportsService.getReportNodeIncidents(id, pageable));
    }

    @GetMapping("/application-incidents/{id}")
    public ResponseEntity<ApplicationIncidentSimplifiedDTO> getApplicationIncidentById(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getApplicationIncidentById(id));
    }

    @GetMapping("/application-incidents/{id}/sources")
    public ResponseEntity<ReportPaginatedIncidentsDTO<ApplicationIncidentSource>> getApplicationIncidentSourcesByIncidentId(
            @PathVariable String id, Pageable pageable) {
        return ResponseEntity.ok(reportsService.getApplicationIncidentSourcesByIncidentId(id, pageable));
    }

    @GetMapping("/node-incidents/{id}")
    public ResponseEntity<NodeIncidentSimplifiedDTO> getNodeIncidentById(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getNodeIncidentById(id));
    }

    @GetMapping("/node-incidents/{id}/sources")
    public ResponseEntity<ReportPaginatedIncidentsDTO<NodeIncidentSource>> getNodeIncidentSourcesByIncidentId(
            @PathVariable String id, Pageable pageable) {
        return ResponseEntity.ok(reportsService.getNodeIncidentSourcesByIncidentId(id, pageable));
    }
}
