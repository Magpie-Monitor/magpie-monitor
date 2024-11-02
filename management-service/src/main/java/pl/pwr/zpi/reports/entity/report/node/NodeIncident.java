package pl.pwr.zpi.reports.entity.report.node;

import lombok.Data;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Data
public class NodeIncident {
    private String id;
    private String category;
    private String clusterId;
    private String nodeName;
    private String summary;
    private String recommendation;
    private Urgency urgency;
    private List<NodeIncidentSource> sources;
}
