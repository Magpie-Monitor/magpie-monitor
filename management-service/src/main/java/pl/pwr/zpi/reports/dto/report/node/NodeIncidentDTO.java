package pl.pwr.zpi.reports.dto.report.node;

import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

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
}
