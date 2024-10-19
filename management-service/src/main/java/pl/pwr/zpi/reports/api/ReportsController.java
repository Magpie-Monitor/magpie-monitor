package pl.pwr.zpi.reports.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import pl.pwr.zpi.reports.ReportDetailedSummary;
import pl.pwr.zpi.reports.dto.report.Report;
import pl.pwr.zpi.reports.dto.report.ReportSummary;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.dto.report.node.NodeIncident;
import pl.pwr.zpi.reports.dto.report.node.ReportIncidents;
import pl.pwr.zpi.reports.service.ReportsService;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/reports")
public class ReportsController {

    private final ReportsService reportsService;

    @GetMapping
    public ResponseEntity<List<ReportSummary>> getReports() {
        return ResponseEntity.ok().body(reportsService.getReportSummaries());
    }

    @GetMapping("/{id}")
    public ResponseEntity<ReportDetailedSummary> getReportById(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getReportDetailedSummaryById(id));
    }

    @GetMapping("/{id}/incidents")
    public ResponseEntity<ReportIncidents> getReportIncidents(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getReportIncidents(id));
    }

    @GetMapping("/application-incidents/{id}")
    public ResponseEntity<List<ApplicationIncident>> getApplicationIncidentById(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getApplicationIncidentById(id));
    }

    @GetMapping("/node-incidents/{id}")
    public ResponseEntity<List<NodeIncident>> getNodeIncidentById(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getNodeIncidentById(id));
    }

}
