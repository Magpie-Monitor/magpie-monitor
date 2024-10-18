package pl.pwr.zpi.reports.dto.report;

import lombok.extern.jackson.Jacksonized;
import pl.pwr.zpi.reports.Urgency;
import pl.pwr.zpi.reports.dto.report.application.ApplicationReport;
import pl.pwr.zpi.reports.dto.report.application.ScheduledApplicationInsight;
import pl.pwr.zpi.reports.dto.report.node.NodeReport;
import pl.pwr.zpi.reports.dto.report.node.ScheduledNodeInsights;

import java.util.List;

public record Report(
        String id,
        String status,
        String clusterId,
        Long sinceMs,
        Long toMs,
        Long requestedAtMs,
        Long scheduledGenerationAtMs,
        String title,
        List<NodeReport> nodeReports,
        List<ApplicationReport> applicationReports,
        Integer totalApplicationEntries,
        Integer totalNodeEntries,
        Urgency urgency,
        List<ScheduledApplicationInsight> scheduledApplicationInsights,
        List<ScheduledNodeInsights> scheduledNodeInsights
) {
}
