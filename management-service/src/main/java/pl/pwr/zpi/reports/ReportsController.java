package pl.pwr.zpi.reports;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import pl.pwr.zpi.reports.scheduler.ReportSchedule;
import pl.pwr.zpi.reports.scheduler.ReportScheduler;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/reports")
public class ReportsController {

    private final ReportsService reportsService;
    private final ReportScheduler reportScheduler;

    @GetMapping
    public ResponseEntity<List<ReportSummarizedDTO>> getReport() throws Exception {
        return ResponseEntity.ok().body(reportsService.getReport());
    }

    @GetMapping("/{id}")
    public ResponseEntity<ReportDTO> getReportById(@PathVariable String id) throws Exception {
        return ResponseEntity.ok().body(reportsService.getReportById(id));
    }

    @PostMapping("/schedule")
    public ResponseEntity<?> setReportSchedule(@RequestParam Long days, @RequestParam Long hours, @RequestParam Long minutes) {
        reportScheduler.scheduleReport(days, hours, minutes);
        return ResponseEntity.ok().build();
    }

    @PostMapping("/schedule/stop")
    public ResponseEntity<?> stopScheduledReport() {
        reportScheduler.deactivateSchedule();
        return ResponseEntity.ok().build();
    }

    @GetMapping("/schedule/last-run")
    public ResponseEntity<ReportSchedule> getLastRunTime() {
        return ResponseEntity.ok().body(reportScheduler.getLastRunTime());
    }
}
