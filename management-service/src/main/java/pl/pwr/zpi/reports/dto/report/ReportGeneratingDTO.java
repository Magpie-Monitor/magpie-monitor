package pl.pwr.zpi.reports.dto.report;

import lombok.Builder;
import pl.pwr.zpi.reports.entity.report.request.ReportGenerationRequestMetadata;
import pl.pwr.zpi.reports.enums.ReportType;

@Builder
public record ReportGeneratingDTO(
        String clusterId,
        ReportType reportType,
        Long sinceMs,
        Long toMs,
        Long requestedAtMs

) {

    public static ReportGeneratingDTO ofReportGenerationRequestMetadata(
            ReportGenerationRequestMetadata reportGenerationRequestMetadata) {
        return ReportGeneratingDTO.builder()
                .clusterId(reportGenerationRequestMetadata.getClusterId())
                .reportType(reportGenerationRequestMetadata.getReportType())
                .sinceMs(reportGenerationRequestMetadata.getCreateReportRequest().sinceMs())
                .toMs(reportGenerationRequestMetadata.getCreateReportRequest().toMs())
                .requestedAtMs(reportGenerationRequestMetadata.getRequestedAt())
                .build();
    }
}
