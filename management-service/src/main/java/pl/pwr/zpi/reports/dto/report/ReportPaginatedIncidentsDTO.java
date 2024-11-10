package pl.pwr.zpi.reports.dto.report;

import lombok.Builder;

import java.util.List;

@Builder
public record ReportPaginatedIncidentsDTO<T>(
        List<T> data,
        Long totalEntries
) {
}
