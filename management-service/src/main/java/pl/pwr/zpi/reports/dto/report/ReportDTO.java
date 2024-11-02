package pl.pwr.zpi.reports.dto.report;

import pl.pwr.zpi.reports.dto.report.application.ApplicationReportDTO;
import pl.pwr.zpi.reports.dto.report.application.ScheduledApplicationInsightDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeReportDTO;
import pl.pwr.zpi.reports.dto.report.node.ScheduledNodeInsightsDTO;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

public record ReportDTO(
        String id,
        String status,
        String clusterId,
        Long sinceMs,
        Long toMs,
        Long requestedAtMs,
        Long scheduledGenerationAtMs,
        String title,
        List<NodeReportDTO> nodeReports,
        List<ApplicationReportDTO> applicationReports,
        Integer totalApplicationEntries,
        Integer totalNodeEntries,
        Urgency urgency,
        List<ScheduledApplicationInsightDTO> scheduledApplicationInsights,
        List<ScheduledNodeInsightsDTO> scheduledNodeInsights
) {
}
