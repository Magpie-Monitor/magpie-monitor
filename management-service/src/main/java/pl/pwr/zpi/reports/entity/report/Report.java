package pl.pwr.zpi.reports.entity.report;

import jakarta.persistence.Id;
import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.application.ApplicationReport;
import pl.pwr.zpi.reports.entity.report.application.scheduled.ScheduledApplicationInsight;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeReport;
import pl.pwr.zpi.reports.entity.report.node.scheduled.ScheduledNodeInsight;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Data
@Builder
public class Report {
    @Id
    private String id;
    private String clusterId;
    private Long sinceMs;
    private Long toMs;
    private Long requestedAtMs;
    private Long scheduledGenerationAtMs;
    private String title;

    private List<NodeReport> nodeReports;
    private List<ApplicationReport> applicationReports;
    private Integer totalApplicationEntries;
    private Integer totalNodeEntries;
    private Urgency urgency;

    private List<ScheduledApplicationInsight> scheduledApplicationInsights;
    private List<ScheduledNodeInsight> scheduledNodeInsights;

    public List<NodeIncident> getNodeIncidents() {
        return nodeReports.stream()
                .map(NodeReport::getNodeIncidents)
                .flatMap(List::stream)
                .toList();
    }

    public List<ApplicationIncident> getApplicationIncidents() {
        return applicationReports.stream()
                .map(ApplicationReport::getIncidents)
                .flatMap(List::stream)
                .toList();
    }
}
