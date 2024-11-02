package pl.pwr.zpi.reports.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import pl.pwr.zpi.reports.service.ReportsService;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;

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
}
