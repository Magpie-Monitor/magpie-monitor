package pl.pwr.zpi.reports.dto.report.node;

public record NodeIncidentSource(
        Long timestamp,
        String content,
        String filename
) {
}
