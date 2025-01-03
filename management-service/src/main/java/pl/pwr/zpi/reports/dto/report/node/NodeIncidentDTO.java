package pl.pwr.zpi.reports.dto.report.node;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Builder
public record NodeIncidentDTO(
        String id,
        String category,
        String clusterId,
        String nodeName,
        String title,
        String summary,
        Accuracy accuracy,
        String customPrompt,
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
                .title(incident.getTitle())
                .summary(incident.getSummary())
                .accuracy(incident.getAccuracy())
                .customPrompt(incident.getCustomPrompt())
                .recommendation(incident.getRecommendation())
                .urgency(incident.getUrgency())
                .sources(incident.getSources().stream()
                        .map(NodeIncidentSourceDTO::fromNodeIncidentSource)
                        .toList()
                )
                .build();
    }
}
