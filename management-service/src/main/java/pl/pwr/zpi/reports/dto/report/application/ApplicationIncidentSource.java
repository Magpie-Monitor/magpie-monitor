package pl.pwr.zpi.reports.dto.report.application;

public record ApplicationIncidentSource(
        Long timestamp,
        String podName,
        String containerName,
        String image,
        String content
) {
}
