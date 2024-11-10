package pl.pwr.zpi.reports.dto.event;

import pl.pwr.zpi.reports.entity.report.Report;

public record ReportGenerated(
        String correlationId,
        Report report,
        Long timestampMs
) {
    public String getReportId() {
        return report.getId();
    }
}
