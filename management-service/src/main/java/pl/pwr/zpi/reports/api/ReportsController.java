package pl.pwr.zpi.reports.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import pl.pwr.zpi.reports.dto.report.Report;
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
    public ResponseEntity<List<Report>> getReport() {
        return ResponseEntity.ok().body(reportsService.getReports());
    }

    @GetMapping("/{id}")
    public ResponseEntity<Report> getReportById(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getReportById(id));
    }

    @GetMapping("/{id}/incidents")
    public ResponseEntity<ReportIncidents> getReportIncidents(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getReportIncidents(id));
    }

    @GetMapping("/application-incidents/{id}")
    public ResponseEntity<List<ApplicationIncident>> getReportApplicationIncidents(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getReportApplicationIncidents(id));
    }

    @GetMapping("/node-incidents/{id}")
    public ResponseEntity<List<NodeIncident>> getReportNodeIncidents(@PathVariable String id) {
        return ResponseEntity.ok().body(reportsService.getReportNodeIncidents(id));
    }


}
