package pl.pwr.zpi.reports.service;

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.dto.report.Report;
import pl.pwr.zpi.reports.dto.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.dto.report.node.NodeIncident;
import pl.pwr.zpi.reports.dto.report.node.ReportIncidents;
import pl.pwr.zpi.utils.client.HttpClient;

import java.util.List;
import java.util.Map;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReportsService {

    @Value("${reports.service.base.url}")
    private String REPORT_SERVICE_BASE_URL;
    private final HttpClient httpClient;

    public List<Report> getReports() {
        String url = String.format("%s/v1/reports", REPORT_SERVICE_BASE_URL);
        return httpClient.getList(
                url,
                Map.of(),
                new TypeReference<>() {
                }
        );
    }

    public Report getReportById(String reportId) {
        String url = String.format("%s/v1/reports/%s", REPORT_SERVICE_BASE_URL, reportId);
        return httpClient.get(
                url,
                Map.of(),
                Report.class
        );
    }

    public List<ApplicationIncident> getReportApplicationIncidents(String reportId) {
        String url = String.format("%s/v1/application-incidents/%s", REPORT_SERVICE_BASE_URL, reportId);
        return httpClient.getList(
                url,
                Map.of(),
                new TypeReference<>() {
                }
        );
    }

    public List<NodeIncident> getReportNodeIncidents(String reportId) {
        String url = String.format("%s/v1/node-incidents/%s", REPORT_SERVICE_BASE_URL, reportId);
        return httpClient.getList(
                url,
                Map.of(),
                new TypeReference<>() {
                }
        );
    }

    public ReportIncidents getReportIncidents(String id) {
        List<NodeIncident> nodeIncidents = getReportNodeIncidents(id);
        List<ApplicationIncident> applicationIncidents = getReportApplicationIncidents(id);
        return new ReportIncidents(applicationIncidents, nodeIncidents);
    }
}
