package pl.pwr.zpi.reports.dto.report.node;

public record NodeReport(
        String node,
        String precision,
        String customPrompt
) {}
