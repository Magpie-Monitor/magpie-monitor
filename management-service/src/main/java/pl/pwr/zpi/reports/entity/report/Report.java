package pl.pwr.zpi.reports.entity.report;

import jakarta.persistence.Id;
import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.dto.report.ReportDTO;
import pl.pwr.zpi.reports.dto.report.application.ApplicationReportDTO;
import pl.pwr.zpi.reports.dto.report.application.ScheduledApplicationInsightDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeReportDTO;
import pl.pwr.zpi.reports.dto.report.node.ScheduledNodeInsightsDTO;
import pl.pwr.zpi.reports.enums.ReportGenerationStatus;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Data
@Builder
public class Report {
    // Used to determine report generation status
    @Id
    private String correlationId;
    private ReportGenerationStatus status;

    private String id;
    private String clusterId;
    private Long sinceMs;
    private Long toMs;
    private Long requestedAtMs;
    private Long scheduledGenerationAtMs;
    private String title;

    private List<NodeReportDTO> nodeReports;
    private List<ApplicationReportDTO> applicationReports;
    private Integer totalApplicationEntries;
    private Integer totalNodeEntries;
    private Urgency urgency;

    private List<ScheduledApplicationInsightDTO> scheduledApplicationInsights;
    private List<ScheduledNodeInsightsDTO> scheduledNodeInsights;

    public static Report generatingReport(String correlationId) {
        return Report.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATING)
                .build();
    }

    public static Report generatedReport(String correlationId, ReportDTO reportDTO) {
        return Report.builder()
                .correlationId(correlationId)
                .status(ReportGenerationStatus.GENERATED)
                .id(reportDTO.id())
                .clusterId(reportDTO.clusterId())
                .sinceMs(reportDTO.sinceMs())
                .toMs(reportDTO.toMs())
                .requestedAtMs(reportDTO.requestedAtMs())
                .scheduledGenerationAtMs(reportDTO.scheduledGenerationAtMs())
                .title(reportDTO.title())
                .nodeReports(reportDTO.nodeReports())
                .applicationReports(reportDTO.applicationReports())
                .totalApplicationEntries(reportDTO.totalApplicationEntries())
                .totalNodeEntries(reportDTO.totalNodeEntries())
                .urgency(reportDTO.urgency())
                .scheduledApplicationInsights(reportDTO.scheduledApplicationInsights())
                .scheduledNodeInsights(reportDTO.scheduledNodeInsights())
                .build();
    }

    public void markAsGenerated() {
        this.status = ReportGenerationStatus.GENERATED;
    }
}
