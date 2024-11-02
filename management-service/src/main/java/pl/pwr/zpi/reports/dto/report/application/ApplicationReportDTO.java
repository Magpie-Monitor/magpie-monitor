package pl.pwr.zpi.reports.dto.report.application;

import java.util.List;

public record ApplicationReportDTO(
        String applicationName,
        String precision,
        String customPrompt,
        List<ApplicationIncidentDTO> incidents
) {
}
