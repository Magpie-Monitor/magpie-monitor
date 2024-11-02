package pl.pwr.zpi.reports.dto.event;

public record ReportRequestFailed(
        String correlationId,
        String errorType,
        String errorMessage,
        Long TimestampMs
) {
}
