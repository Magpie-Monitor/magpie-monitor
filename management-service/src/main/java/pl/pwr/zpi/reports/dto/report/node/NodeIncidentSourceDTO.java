package pl.pwr.zpi.reports.dto.report.node;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource;

@Builder
public record NodeIncidentSourceDTO(
        Long timestamp,
        String content,
        String filename
) {
    public static NodeIncidentSourceDTO fromNodeIncidentSource(NodeIncidentSource incidentSource) {
        return NodeIncidentSourceDTO.builder()
                .timestamp(incidentSource.getTimestamp())
                .content(incidentSource.getContent())
                .filename(incidentSource.getFilename())
                .build();
    }
}
