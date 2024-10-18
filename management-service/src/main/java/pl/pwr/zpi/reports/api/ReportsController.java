package pl.pwr.zpi.reports.api;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import pl.pwr.zpi.reports.dto.ReportSummarized;
import pl.pwr.zpi.reports.service.ReportsService;
import pl.pwr.zpi.reports.dto.Report;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/reports")
public class ReportsController {

    private final ReportsService reportsService;

    @GetMapping
    public ResponseEntity<List<ReportSummarized>> getReport() throws Exception {
        return ResponseEntity.ok().body(reportsService.getReport());
    }

    @GetMapping("/{id}")
    public ResponseEntity<Report> getReportById(@PathVariable String id) throws Exception {
        return ResponseEntity.ok().body(reportsService.getReportById(id));
    }
}
