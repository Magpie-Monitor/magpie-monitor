package pl.pwr.zpi.reports.entity.report.node;

import jakarta.persistence.Id;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NodeIncident {
    @Id
    private String id;
    private String title;
    private Accuracy accuracy;
    private String customPrompt;
    private String category;
    private String clusterId;
    private String nodeName;
    private String summary;
    private String recommendation;
    private Urgency urgency;
    private List<NodeIncidentSource> sources;
}
