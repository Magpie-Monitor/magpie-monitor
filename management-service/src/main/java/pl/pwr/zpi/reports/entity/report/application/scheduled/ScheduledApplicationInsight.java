package pl.pwr.zpi.reports.entity.report.application.scheduled;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;
import java.util.Map;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ScheduledApplicationInsight {
    private List<String> scheduledJobIds;
    private Long sinceMs;
    private Long toMs;
    private String clusterId;
    private Map<String, ScheduledApplicationInsightConfiguration> applicationConfiguration;
//    private List<ScheduledApplicationInsightConfiguration> applicationConfiguration;
}
