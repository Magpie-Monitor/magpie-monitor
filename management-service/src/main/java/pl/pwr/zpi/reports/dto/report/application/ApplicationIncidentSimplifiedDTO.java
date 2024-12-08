package pl.pwr.zpi.reports.dto.report.application;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

@Builder
public record ApplicationIncidentSimplifiedDTO(
        String id,
        String reportId,
        String title,
        Accuracy accuracy,
        String customPrompt,
        String clusterId,
        String applicationName,
        String category,
        String summary,
        String recommendation,
        Urgency urgency,
        long sinceMs,
        long toMs
) {
    public static ApplicationIncidentSimplifiedDTO fromApplicationIncident(ApplicationIncident incident) {
        return ApplicationIncidentSimplifiedDTO.builder()
                .id(incident.getId())
                .reportId(incident.getReportId())
                .title(incident.getTitle())
                .accuracy(incident.getAccuracy())
                .customPrompt(incident.getCustomPrompt())
                .clusterId(incident.getClusterId())
                .applicationName(incident.getApplicationName())
                .category(incident.getCategory())
                .summary(incident.getSummary())
                .recommendation(incident.getRecommendation())
                .urgency(incident.getUrgency())
                .sinceMs(incident.getSources().stream()
                        .map(ApplicationIncidentSource::getTimestamp)
                        .min(Long::compareTo)
                        .orElse(0L))
                .toMs(incident.getSources().stream()
                        .map(ApplicationIncidentSource::getTimestamp)
                        .max(Long::compareTo)
                        .orElse(0L))
                .build();
    }
}