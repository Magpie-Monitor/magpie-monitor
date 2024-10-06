package pl.pwr.zpi.reports;

public record ReportSummarizedDTO(String id, String cluster, String title, String summary, Urgency urgency, Long fromDateNs, Long toDateNs) {
}
