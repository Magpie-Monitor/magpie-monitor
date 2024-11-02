package pl.pwr.zpi.reports.dto.report.node;

import java.util.List;

public record ScheduledNodeInsightsDTO(
        List<String> scheduledJobIds,
        Long sinceMs,
        Long toMs,
        String clusterId,
        List<ScheduledNodeInsightConfigurationDTO> nodeConfiguration
) {
}
