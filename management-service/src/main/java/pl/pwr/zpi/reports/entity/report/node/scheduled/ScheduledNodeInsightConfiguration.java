package pl.pwr.zpi.reports.entity.report.node.scheduled;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.reports.enums.Accuracy;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ScheduledNodeInsightConfiguration {
    private String nodeName;
    private Accuracy accuracy;
    private String customPrompt;
}
