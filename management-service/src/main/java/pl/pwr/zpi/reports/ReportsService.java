package pl.pwr.zpi.reports;

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class ReportsService {
    private final ReportsClient reportsClient;
    public List<ReportSummarizedDTO> getReport() throws Exception {
        log.info("Getting report");
        return reportsClient.sendGetRequestForList(
                "http://reports-service:8099/v1/reports",
                new TypeReference<>() {}
        );
    }

    public ReportDTO getReportById(String id) throws Exception {
        log.info("Getting report by id: {}", id);
        return reportsClient.sendGetRequest(
                "http://reports-service:8099/v1/reports/" + id,
                ReportDTO.class
        );      }
}
