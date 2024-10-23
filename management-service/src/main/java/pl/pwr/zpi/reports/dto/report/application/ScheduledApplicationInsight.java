package pl.pwr.zpi.reports.dto.report.application;

import java.util.List;

public record ScheduledApplicationInsight(
        List<String> scheduledJobIds,
        Long sinceMs,
        Long toMs,
        String clusterId,
        List<ScheduledApplicationInsightConfiguration> applicationConfiguration
) {
}
