package pl.pwr.zpi.reports.dto.report.node;

public record NodeIncidentSourceDTO(
        Long timestamp,
        String content,
        String filename
) {
}
