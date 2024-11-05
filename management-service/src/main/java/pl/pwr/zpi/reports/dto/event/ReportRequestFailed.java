package pl.pwr.zpi.reports.dto.event;

import lombok.Builder;

@Builder
public record ReportRequestFailed(
        String correlationId,
        String errorType,
        String errorMessage,
        Long timestampMs
) {
}
