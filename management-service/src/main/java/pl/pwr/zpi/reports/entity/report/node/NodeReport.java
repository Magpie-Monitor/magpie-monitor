package pl.pwr.zpi.reports.entity.report.node;

import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.enums.Precision;

import java.util.List;

@Data
@Builder
public class NodeReport {
    private String node;
    private Precision precision;
    private String customPrompt;
    private List<NodeIncident> incidents;
}
