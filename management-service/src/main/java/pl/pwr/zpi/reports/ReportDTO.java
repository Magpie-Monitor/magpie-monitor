package pl.pwr.zpi.reports;

public record ReportDTO(String id, String cluster, String title, String summary, Urgency urgency, Long fromDateNs,
                        Long toDateNs, Long totalApplicationEntries, Long totalNodeEntries) {
}
//  TBA: Statistics statistics