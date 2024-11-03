package pl.pwr.zpi.reports.dto.report.node;

import pl.pwr.zpi.reports.dto.report.Accuracy;
import pl.pwr.zpi.reports.dto.report.Urgency;

import java.util.List;

public record NodeIncident(
        String id,
        String title,
        Accuracy accuracy,
        String customPrompt,
        String category,
        String clusterId,
        String nodeName,
        String summary,
        String recommendation,
        Urgency urgency,
        List<NodeIncidentSource> sources
) {
}
