package pl.pwr.zpi.reports.dto.request;

import jakarta.validation.constraints.Min;
import lombok.Builder;
import lombok.NonNull;

@Builder
public record CreateReportScheduleRequest(
        @NonNull
        String clusterId,
        @NonNull
        @Min(value = 86400000, message = "Period must be at least 24 hours.")
        Long periodMs
) {}
