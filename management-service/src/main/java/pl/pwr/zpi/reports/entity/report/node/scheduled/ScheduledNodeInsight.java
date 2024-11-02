package pl.pwr.zpi.reports.entity.report.node.scheduled;

import lombok.Data;

import java.util.List;

@Data
public class ScheduledNodeInsight {
    private List<String> scheduledJobIds;
    private Long sinceMs;
    private Long toMs;
    private String clusterId;
    private List<ScheduledNodeInsightConfiguration> nodeConfiguration;
}
