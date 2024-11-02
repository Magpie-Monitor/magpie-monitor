package pl.pwr.zpi.reports.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.reports.dto.report.ReportDetailedSummaryDTO;
import pl.pwr.zpi.reports.dto.report.ReportIncidentsDTO;
import pl.pwr.zpi.reports.dto.report.ReportSummaryDTO;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.service.ReportsService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/reports")
public class ReportsController {

    private final ReportsService reportsService;

//    @GetMapping
//    public ResponseEntity<List<ReportSummary>> getReports() {
//        return ResponseEntity.ok().body(reportsService.getReportSummaries());
//    }

//    @GetMapping("/{id}")
//    public ResponseEntity<ReportDetailedSummary> getReportById(@PathVariable String id) {
//        return ResponseEntity.ok().body(reportsService.getReportDetailedSummaryById(id));
//    }

//    @GetMapping("/{id}/incidents")
//    public ResponseEntity<ReportIncidents> getReportIncidents(@PathVariable String id) {
//        return ResponseEntity.ok().body(reportsService.getReportIncidents(id));
//    }

    @GetMapping("/application-incidents/{id}")
    public ResponseEntity<ApplicationIncident> getApplicationIncidentById(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getApplicationIncidentById(id));
    }

    @GetMapping("/node-incidents/{id}")
    public ResponseEntity<NodeIncident> getNodeIncidentById(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getNodeIncidentById(id));
    }

    @PostMapping
    public ResponseEntity<Void> createReport(@RequestBody CreateReportRequest reportRequest) {
        reportsService.createReport(reportRequest);
        return ResponseEntity.ok().build();
    }

    @GetMapping
    public ResponseEntity<List<ReportSummaryDTO>> getReports() {
        return ResponseEntity.ok().body(reportsService.getReportSummaries());
    }

    @GetMapping("/{id}")
    public ResponseEntity<ReportDetailedSummaryDTO> getReportById(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getReportDetailedSummaryById(id));
    }

    @GetMapping("/{id}/incidents")
    public ResponseEntity<ReportIncidentsDTO> getReportIncidents(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getReportIncidents(id));
    }

    @GetMapping("/application-incidents/{id}")
    public ResponseEntity<ApplicationIncident> getApplicationIncidentById(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getApplicationIncidentById(id));
    }

    @GetMapping("/node-incidents/{id}")
    public ResponseEntity<NodeIncident> getNodeIncidentById(@PathVariable String id) {
        return ResponseEntity.of(reportsService.getNodeIncidentById(id));
    }
}
