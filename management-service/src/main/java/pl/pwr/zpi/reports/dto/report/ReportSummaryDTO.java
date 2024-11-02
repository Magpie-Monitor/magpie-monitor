package pl.pwr.zpi.reports.dto.report;

import pl.pwr.zpi.reports.enums.Urgency;

public record ReportSummaryDTO(
        String id,
        String clusterId,
        String title,
        Urgency urgency,
        Long sinceMs,
        Long toMs) {
}
