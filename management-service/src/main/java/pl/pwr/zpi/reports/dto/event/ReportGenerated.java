package pl.pwr.zpi.reports.dto.event;

import pl.pwr.zpi.reports.dto.report.ReportDTO;

public record ReportGenerated(
        String correlationId,
        ReportDTO report,
        Long timestampMs
) {
}
