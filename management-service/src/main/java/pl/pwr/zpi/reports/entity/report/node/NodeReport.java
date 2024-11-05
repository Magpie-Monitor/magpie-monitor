package pl.pwr.zpi.reports.entity.report.node;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NodeReport {
    private String node;
    private Accuracy accuracy;
    private String customPrompt;
    private List<NodeIncident> incidents;
}
