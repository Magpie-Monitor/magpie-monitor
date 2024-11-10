package pl.pwr.zpi.notifications;

public interface ReportNotifier {
    void notifyOnReportGenerated(Long receiverId, String reportId);

    void notifyOnReportGenerationFailed(Long receiverId, String clusterId);
}