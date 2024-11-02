package pl.pwr.zpi.reports.dto.report.application;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource;

@Builder
public record ApplicationIncidentSourceDTO(
        Long timestamp,
        String podName,
        String containerName,
        String image,
        String content
) {
    public static ApplicationIncidentSourceDTO fromApplicationIncidentSource(ApplicationIncidentSource incidentSource) {
        return ApplicationIncidentSourceDTO.builder()
                .timestamp(incidentSource.getTimestamp())
                .podName(incidentSource.getPodName())
                .containerName(incidentSource.getContainerName())
                .image(incidentSource.getImage())
                .content(incidentSource.getContent())
                .build();
    }
}
