package pl.pwr.zpi.reports.repository.projection;

import pl.pwr.zpi.reports.enums.Urgency;

public interface ReportDetailedSummaryProjection {
    String getId();

    String getClusterId();

    String getTitle();

    Urgency getUrgency();

    Long getSinceMs();

    Long getToMs();

    Integer getTotalApplicationEntries();

    Integer getTotalNodeEntries();
}
