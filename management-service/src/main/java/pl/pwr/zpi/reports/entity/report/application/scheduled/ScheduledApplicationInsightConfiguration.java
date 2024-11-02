package pl.pwr.zpi.reports.entity.report.application.scheduled;

import lombok.Data;

@Data
public class ScheduledApplicationInsightConfiguration {
    private String applicationName;
    private String precision;
    private String customPrompt;
}
