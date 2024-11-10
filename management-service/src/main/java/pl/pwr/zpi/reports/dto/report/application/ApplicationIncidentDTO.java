package pl.pwr.zpi.reports.dto.report.application;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Builder
public record ApplicationIncidentDTO(
        String id,
        String clusterId,
        String applicationName,
        String category,
        String title,
        String summary,
        Accuracy accuracy,
        String customPrompt,
        String recommendation,
        Urgency urgency,
        List<ApplicationIncidentSourceDTO> sources
) {
    public static ApplicationIncidentDTO fromApplicationIncident(ApplicationIncident incident) {
        return ApplicationIncidentDTO.builder()
                .id(incident.getId())
                .clusterId(incident.getClusterId())
                .applicationName(incident.getApplicationName())
                .category(incident.getCategory())
                .title(incident.getTitle())
                .summary(incident.getSummary())
                .accuracy(incident.getAccuracy())
                .customPrompt(incident.getCustomPrompt())
                .recommendation(incident.getRecommendation())
                .urgency(incident.getUrgency())
                .sources(incident.getSources()
                        .stream()
                        .map(ApplicationIncidentSourceDTO::fromApplicationIncidentSource)
                        .toList()
                )
                .build();
    }
}
