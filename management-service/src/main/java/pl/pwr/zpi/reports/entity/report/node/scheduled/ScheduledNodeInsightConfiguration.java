package pl.pwr.zpi.reports.entity.report.node.scheduled;

import lombok.Data;

@Data
public class ScheduledNodeInsightConfiguration {
    private String nodeName;
    private String precision;
    private String customPrompt;
}
