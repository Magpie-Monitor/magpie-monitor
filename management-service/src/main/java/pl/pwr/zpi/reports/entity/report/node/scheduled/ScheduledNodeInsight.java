package pl.pwr.zpi.reports.entity.report.node.scheduled;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ScheduledNodeInsight {
    private List<String> scheduledJobIds;
    private Long sinceMs;
    private Long toMs;
    private String clusterId;
    private List<ScheduledNodeInsightConfiguration> nodeConfiguration;
}
