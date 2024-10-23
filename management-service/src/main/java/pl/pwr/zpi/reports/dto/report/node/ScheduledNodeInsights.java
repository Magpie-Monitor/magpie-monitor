package pl.pwr.zpi.reports.dto.report.node;

import java.util.List;

public record ScheduledNodeInsights(
        List<String> scheduledJobIds,
        Long sinceMs,
        Long toMs,
        String clusterId,
        List<ScheduledNodeInsightConfiguration> nodeConfiguration
) {
}
