package pl.pwr.zpi.reports.dto.report;

import lombok.Builder;
import pl.pwr.zpi.reports.enums.Urgency;
import pl.pwr.zpi.reports.repository.projection.ReportSummaryProjection;

@Builder
public record ReportSummaryDTO(
        String id,
        String clusterId,
        String title,
        Urgency urgency,
        Long sinceMs,
        Long toMs) {

    public static ReportSummaryDTO ofReportSummaryProjection(ReportSummaryProjection reportSummaryProjection) {
        return ReportSummaryDTO.builder()
                .id(reportSummaryProjection.getId())
                .clusterId(reportSummaryProjection.getClusterId())
                .title(reportSummaryProjection.getTitle())
                .urgency(reportSummaryProjection.getUrgency())
                .sinceMs(reportSummaryProjection.getSinceMs())
                .toMs(reportSummaryProjection.getToMs())
                .build();
    }
}
