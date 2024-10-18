package pl.pwr.zpi.reports.dto.report.node;

import java.util.List;

public record NodeReport(
        String node,
        String precision,
        String customPrompt,
        List<NodeIncident> nodeIncidents
) {}
