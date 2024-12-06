package pl.pwr.zpi.reports.dto.report.node;

import lombok.Builder;
import pl.pwr.zpi.reports.dto.report.ReportDetailedSummaryDTO;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource;
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource;
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;

import java.util.List;

@Builder
public record ReportDetailedWithIncidentsDTO(
        ReportDetailedSummaryDTO reportDetailedSummary,
        List<ApplicationIncidentSimplifiedDTO> applicationIncidents,
        List<NodeIncidentSimplifiedDTO> nodeIncidents
) {
    @Builder
    record ApplicationIncidentSimplifiedDTO(
            String id,
            String reportId,
            String title,
            Accuracy accuracy,
            String customPrompt,
            String clusterId,
            String applicationName,
            String category,
            String summary,
            String recommendation,
            Urgency urgency,
            long sinceMs,
            long toMs
    ) {
    }

    @Builder
    record NodeIncidentSimplifiedDTO(
            String id,
            String reportId,
            String title,
            Accuracy accuracy,
            String customPrompt,
            String clusterId,
            String nodeName,
            String category,
            String summary,
            String recommendation,
            Urgency urgency,
            long sinceMs,
            long toMs
    ) {
    }

    public static ReportDetailedWithIncidentsDTO mapToReportDetailedWithIncidentsDTO(Report report) {
        return ReportDetailedWithIncidentsDTO.builder()
                .reportDetailedSummary(
                        ReportDetailedSummaryDTO.builder()
                                .id(report.getId())
                                .clusterId(report.getClusterId())
                                .title(report.getTitle())
                                .urgency(report.getUrgency())
                                .requestedAtMs(report.getRequestedAtMs())
                                .sinceMs(report.getSinceMs())
                                .toMs(report.getToMs())
                                .totalApplicationEntries(report.getTotalApplicationEntries())
                                .totalNodeEntries(report.getTotalNodeEntries())
                                .analyzedApplications(report.getAnalyzedApplications())
                                .analyzedNodes(report.getAnalyzedNodes())
                                .build()
                )
                .applicationIncidents(report.getApplicationIncidents().stream()
                        .map(appIncident -> ApplicationIncidentSimplifiedDTO.builder()
                                .id(appIncident.getId())
                                .reportId(appIncident.getReportId())
                                .title(appIncident.getTitle())
                                .accuracy(appIncident.getAccuracy())
                                .customPrompt(appIncident.getCustomPrompt())
                                .clusterId(appIncident.getClusterId())
                                .applicationName(appIncident.getApplicationName())
                                .category(appIncident.getCategory())
                                .summary(appIncident.getSummary())
                                .recommendation(appIncident.getRecommendation())
                                .urgency(appIncident.getUrgency())
                                .sinceMs(appIncident.getSources().stream()
                                        .map(ApplicationIncidentSource::getTimestamp)
                                        .min(Long::compareTo)
                                        .orElse(report.getSinceMs()))
                                .toMs(appIncident.getSources().stream()
                                        .map(ApplicationIncidentSource::getTimestamp)
                                        .max(Long::compareTo)
                                        .orElse(report.getToMs()))
                                .build())
                        .toList()
                )
                .nodeIncidents(report.getNodeIncidents().stream()
                        .map(nodeIncident -> ReportDetailedWithIncidentsDTO.NodeIncidentSimplifiedDTO.builder()
                                .id(nodeIncident.getId())
                                .reportId(nodeIncident.getReportId())
                                .title(nodeIncident.getTitle())
                                .accuracy(nodeIncident.getAccuracy())
                                .customPrompt(nodeIncident.getCustomPrompt())
                                .clusterId(nodeIncident.getClusterId())
                                .nodeName(nodeIncident.getNodeName())
                                .category(nodeIncident.getCategory())
                                .summary(nodeIncident.getSummary())
                                .recommendation(nodeIncident.getRecommendation())
                                .urgency(nodeIncident.getUrgency())
                                .sinceMs(nodeIncident.getSources().stream()
                                        .map(NodeIncidentSource::getTimestamp)
                                        .min(Long::compareTo)
                                        .orElse(report.getSinceMs()))
                                .toMs(nodeIncident.getSources().stream()
                                        .map(NodeIncidentSource::getTimestamp)
                                        .max(Long::compareTo)
                                        .orElse(report.getToMs()))
                                .build())
                        .toList()
                )
                .build();
    }
}