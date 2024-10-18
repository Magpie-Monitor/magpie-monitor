package pl.pwr.zpi.reports.dto.report.application;

import java.util.List;

public record ScheduledApplicationInsight(
        String id,
        Long sinceMs,
        Long toMs,
        String clusterId,
        List<ScheduledApplicationInsightConfiguration> applicationConfiguration
) {
}
