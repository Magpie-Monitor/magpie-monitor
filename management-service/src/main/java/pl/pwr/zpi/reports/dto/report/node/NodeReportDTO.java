package pl.pwr.zpi.reports.dto.report.node;

import java.util.List;

public record NodeReportDTO(
        String node,
        String precision,
        String customPrompt,
        List<NodeIncidentDTO> nodeIncidents
) {}
