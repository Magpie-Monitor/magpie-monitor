package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.broker.ReportPublisher;
import pl.pwr.zpi.reports.dto.event.ReportGenerated;
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed;
import pl.pwr.zpi.reports.dto.event.ReportRequested;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.repository.ReportRepository;

import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ReportsService {

    private final ReportRepository reportRepository;
    private final ReportPublisher reportPublisher;

    public void createReport(CreateReportRequest reportRequest) {
        ReportRequested reportRequested = ReportRequested.of(reportRequest);
        reportPublisher.publishReportRequestedEvent(reportRequested);

        Report report = Report.generatingReport(reportRequested.correlationId());
        persistReport(report);
    }

    public void handleGenerationError(ReportRequestFailed requestFailed) {
        Optional<Report> report = reportRepository.findByCorrelationId(requestFailed.correlationId());
        report.ifPresent(r -> {
            r.markAsGenerated();
            persistReport(r);
        });
    }

    public void handleReportGenerated(ReportGenerated reportGenerated) {
        Report report = Report.generatedReport(reportGenerated.correlationId(), reportGenerated.report());
        persistReport(report);
    }


    private Report persistReport(Report report) {
        return reportRepository.save(report);
    }

//    private final HttpClient httpClient;

//    public <T> T getReportRepresentationById(String reportId, Class<T> clazz) {
//        String url = String.format("%s/v1/reports/%s", REPORT_SERVICE_BASE_URL, reportId);
//        return httpClient.get(
//                url,
//                Collections.emptyMap(),
//                clazz
//        );
//    }
//
//    public <T> List<T> getReportListRepresentation(TypeReference<List<T>> typeReference) {
//        String url = String.format("%s/v1/reports", REPORT_SERVICE_BASE_URL);
//        return httpClient.getList(
//                url,
//                Collections.emptyMap(),
//                typeReference
//        );
//    }
//
//    public List<ReportSummary> getReportSummaries() {
//        return getReportListRepresentation(new TypeReference<>() {
//        });
//    }
//
//    public ReportDetailedSummary getReportDetailedSummaryById(String reportId) {
//        return getReportRepresentationById(reportId, ReportDetailedSummary.class);
//    }
//
//    public List<ApplicationIncident> getApplicationIncidentById(String incidentId) {
//        String url = String.format("%s/v1/application-incidents/%s", REPORT_SERVICE_BASE_URL, incidentId);
//        return httpClient.getList(
//                url,
//                Collections.emptyMap(),
//                new TypeReference<>() {
//                }
//        );
//    }
//
//    public List<NodeIncident> getNodeIncidentById(String incidentId) {
//        String url = String.format("%s/v1/node-incidents/%s", REPORT_SERVICE_BASE_URL, incidentId);
//        return httpClient.getList(
//                url,
//                Collections.emptyMap(),
//                new TypeReference<>() {
//                }
//        );
//    }
//
//    public ReportIncidents getReportIncidents(String id) {
//        Report report = getReportRepresentationById(id, Report.class);
//        return new ReportIncidents(
//                report.applicationReports().stream().
//                        flatMap(applicationReport -> applicationReport.incidents().stream())
//                        .toList(),
//                report.nodeReports().stream()
//                        .flatMap(nodeReport -> nodeReport.nodeIncidents().stream())
//                        .toList()
//        );
//    }

}
