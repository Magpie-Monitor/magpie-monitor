package pl.pwr.zpi.reports.service;

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.ReportDetailedSummary;
import pl.pwr.zpi.reports.dto.report.Report;
import pl.pwr.zpi.reports.dto.report.ReportSummary;
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

    public <T> T getReportRepresentationById(String reportId, Class<T> clazz) {
        String url = String.format("%s/v1/reports/%s", REPORT_SERVICE_BASE_URL, reportId);
        return httpClient.get(
                url,
                Map.of(),
                clazz
        );
    }

    public <T> List<T> getReportListRepresentation(TypeReference<List<T>> typeReference) {
        String url = String.format("%s/v1/reports", REPORT_SERVICE_BASE_URL);
        return httpClient.getList(
                url,
                Map.of(),
                typeReference
        );
    }

    public List<ReportSummary> getReportSummaries() {
        return getReportListRepresentation(new TypeReference<>() {
        });
    }

    public ReportDetailedSummary getReportDetailedSummaryById(String reportId) {
        return getReportRepresentationById(reportId, ReportDetailedSummary.class);
    }

    public List<ApplicationIncident> getApplicationIncidentById(String incidentId) {
        String url = String.format("%s/v1/application-incidents/%s", REPORT_SERVICE_BASE_URL, incidentId);
        return httpClient.getList(
                url,
                Map.of(),
                new TypeReference<>() {
                }
        );
    }

    public List<NodeIncident> getNodeIncidentById(String incidentId) {
        String url = String.format("%s/v1/node-incidents/%s", REPORT_SERVICE_BASE_URL, incidentId);
        return httpClient.getList(
                url,
                Map.of(),
                new TypeReference<>() {
                }
        );
    }

    public ReportIncidents getReportIncidents(String id) {
        Report report = getReportRepresentationById(id, Report.class);
        return new ReportIncidents(
                report.applicationReports().stream().
                        flatMap(applicationReport -> applicationReport.incidents().stream())
                        .toList(),
                report.nodeReports().stream()
                        .flatMap(nodeReport -> nodeReport.nodeIncidents().stream())
                        .toList()
        );
    }

}
