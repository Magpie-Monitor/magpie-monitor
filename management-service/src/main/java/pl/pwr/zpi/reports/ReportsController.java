package pl.pwr.zpi.reports;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/reports")
public class ReportsController {

    private final ReportsService reportsService;

    @GetMapping
    public ResponseEntity<List<ReportSummarizedDTO>> getReport() throws Exception {
        return ResponseEntity.ok().body(reportsService.getReport());
    }

    @GetMapping("/{id}")
    public ResponseEntity<ReportDTO> getReportById(@PathVariable String id) throws Exception {
        return ResponseEntity.ok().body(reportsService.getReportById(id));
    }
}
