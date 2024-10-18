package pl.pwr.zpi.reports.dto.report.application;

import java.util.List;

public record ApplicationReport(
        String applicationName,
        String precision,
        String customPrompt,
        List<ApplicationIncident> incidents
) {
}
