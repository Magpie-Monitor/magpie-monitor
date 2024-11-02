package pl.pwr.zpi.reports.entity.report.application.scheduled;

import lombok.Data;

import java.util.List;

@Data
public class ScheduledApplicationInsight {
    private List<String> scheduledJobIds;
    private Long sinceMs;
    private Long toMs;
    private String clusterId;
    private List<ScheduledApplicationInsightConfiguration> applicationConfiguration;
}
