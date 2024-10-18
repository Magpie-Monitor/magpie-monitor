package pl.pwr.zpi.reports.dto;

import pl.pwr.zpi.reports.Urgency;

public record Report(String id, String clusterId, String title, String summary, Urgency urgency, Long fromDateNs,
                     Long toDateNs, Long totalApplicationEntries, Long totalNodeEntries) {
}
//  TBA: Statistics statistics