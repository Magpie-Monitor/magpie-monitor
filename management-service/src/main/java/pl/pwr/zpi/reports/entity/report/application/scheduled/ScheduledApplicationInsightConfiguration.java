package pl.pwr.zpi.reports.entity.report.application.scheduled;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.reports.enums.Accuracy;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ScheduledApplicationInsightConfiguration {
    private String applicationName;
    private Accuracy accuracy;
    private String customPrompt;
}
