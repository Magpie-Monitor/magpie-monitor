package pl.pwr.zpi.reports.dto.report;

import lombok.Builder;
import pl.pwr.zpi.reports.enums.Urgency;
import pl.pwr.zpi.reports.repository.projection.ReportDetailedSummaryProjection;

@Builder
public record ReportDetailedSummaryDTO(
        String id,
        String clusterId,
        String title,
        Urgency urgency,
        Long requestedAtMs,
        Long sinceMs,
        Long toMs,
        Integer totalApplicationEntries,
        Integer totalNodeEntries,
        Integer analyzedApplications,
        Integer analyzedNodes
) {

    public static ReportDetailedSummaryDTO fromReportDetailedSummaryProjection(ReportDetailedSummaryProjection projection) {
        return ReportDetailedSummaryDTO.builder()
                .id(projection.getId())
                .clusterId(projection.getClusterId())
                .title(projection.getTitle())
                .urgency(projection.getUrgency())
                .requestedAtMs(projection.getRequestedAtMs())
                .sinceMs(projection.getSinceMs())
                .toMs(projection.getToMs())
                .totalApplicationEntries(projection.getTotalApplicationEntries())
                .totalNodeEntries(projection.getTotalNodeEntries())
                .analyzedApplications(projection.getAnalyzedApplications())
                .analyzedNodes(projection.getAnalyzedNodes())
                .build();
    }
}
