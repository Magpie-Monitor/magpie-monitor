package pl.pwr.zpi.reports;

public record ReportDTO(String id, String clusterId, String title, String summary, Urgency urgency, Long fromDateNs,
                        Long toDateNs, Long totalApplicationEntries, Long totalNodeEntries) {
}
//  TBA: Statistics statistics