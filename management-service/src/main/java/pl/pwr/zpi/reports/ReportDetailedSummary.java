package pl.pwr.zpi.reports;

import pl.pwr.zpi.reports.dto.report.Urgency;

public record ReportDetailedSummary(
        String id,
        String clusterId,
        String title,
        Urgency urgency,
        Long sinceMs,
        Long toMs,
        Integer totalApplicationEntries,
        Integer totalNodeEntries) {
}
// TODO - evaluate whether summary and statistics fields are needed
