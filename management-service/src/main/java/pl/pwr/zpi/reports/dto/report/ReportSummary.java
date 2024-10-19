package pl.pwr.zpi.reports.dto.report;

public record ReportSummary(
        String id,
        String clusterId,
        String title,
        Urgency urgency,
        Long sinceMs,
        Long toMs) {
}
