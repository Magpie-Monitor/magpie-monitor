package pl.pwr.zpi.reports;

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class ReportsService {
    
    private final ReportsClient reportsClient;
    // TODO - refactor
    private final String ALL_REPORTS_URL = "http://reports-service:8099/v1/reports";
    private final String REPORT_DETAILS_URL = "http://reports-service:8099/v1/reports/";

    public List<ReportSummarizedDTO> getReport() throws Exception {
        log.info("Getting report");
        return reportsClient.sendGetRequestForList(
                ALL_REPORTS_URL,
                new TypeReference<>() {}
        );
    }

    public ReportDTO getReportById(String id) throws Exception {
        log.info("Getting report by id: {}", id);
        return reportsClient.sendGetRequest(
                REPORT_DETAILS_URL + id,
                ReportDTO.class
        );      }
}
