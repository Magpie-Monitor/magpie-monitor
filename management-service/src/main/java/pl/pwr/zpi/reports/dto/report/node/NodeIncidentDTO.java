package pl.pwr.zpi.reports.dto.report.node;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Builder
public record NodeIncidentDTO(
        String id,
        String category,
        String clusterId,
        String nodeName,
        String summary,
        String recommendation,
        Urgency urgency,
        List<NodeIncidentSourceDTO> sources
) {
    public static NodeIncidentDTO fromNodeIncident(NodeIncident incident) {
        return NodeIncidentDTO.builder()
                .id(incident.getId())
                .category(incident.getCategory())
                .clusterId(incident.getClusterId())
                .nodeName(incident.getNodeName())
                .summary(incident.getSummary())
                .recommendation(incident.getRecommendation())
                .urgency(incident.getUrgency())
                .sources(incident.getSources().stream()
                        .map(NodeIncidentSourceDTO::fromNodeIncidentSource)
                        .toList()
                )
                .build();
    }
}
