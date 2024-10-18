package pl.pwr.zpi.reports.dto.report.node;

import java.util.List;

public record ScheduledNodeInsights(
        String id,
        Long sinceMs,
        Long toMs,
        String clusterId,
        List<ScheduledNodeInsightConfiguration> nodeConfiguration
) {
}
