package pl.pwr.zpi.reports.dto.report.node;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

@Builder
public record NodeIncidentSimplifiedDTO(
        String id,
        String reportId,
        String title,
        Accuracy accuracy,
        String customPrompt,
        String clusterId,
        String nodeName,
        String category,
        String summary,
        String recommendation,
        Urgency urgency,
        long sinceMs,
        long toMs
) {
    public static NodeIncidentSimplifiedDTO fromNodeIncident(NodeIncident incident) {
        return NodeIncidentSimplifiedDTO.builder()
                .id(incident.getId())
                .reportId(incident.getReportId())
                .title(incident.getTitle())
                .accuracy(incident.getAccuracy())
                .customPrompt(incident.getCustomPrompt())
                .clusterId(incident.getClusterId())
                .nodeName(incident.getNodeName())
                .category(incident.getCategory())
                .summary(incident.getSummary())
                .recommendation(incident.getRecommendation())
                .urgency(incident.getUrgency())
                .sinceMs(incident.getSources().stream()
                        .map(NodeIncidentSource::getTimestamp)
                        .min(Long::compareTo)
                        .orElse(0L))
                .toMs(incident.getSources().stream()
                        .map(NodeIncidentSource::getTimestamp)
                        .max(Long::compareTo)
                        .orElse(0L))
                .build();
    }
}