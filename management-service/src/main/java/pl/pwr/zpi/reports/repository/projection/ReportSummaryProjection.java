package pl.pwr.zpi.reports.repository.projection;

import pl.pwr.zpi.reports.enums.Urgency;

public interface ReportSummaryProjection {
    String getId();

    String getClusterId();

    String getTitle();

    Urgency getUrgency();

    Long getRequestedAtMs();

    Long getSinceMs();

    Long getToMs();
}
