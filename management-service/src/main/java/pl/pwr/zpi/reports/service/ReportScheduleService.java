package pl.pwr.zpi.reports.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.reports.dto.request.CreateReportScheduleRequest;
import pl.pwr.zpi.reports.dto.scheduler.ReportSchedule;
import pl.pwr.zpi.reports.repository.ReportScheduleRepository;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReportScheduleService {
    private final ReportScheduleRepository reportScheduleRepository;
    private final ClusterRepository clusterRepository;

    public void scheduleReport(CreateReportScheduleRequest scheduleRequest) {
        validateClusterId(scheduleRequest.clusterId());
        reportScheduleRepository.save(ReportSchedule.fromCreateScheduleRequest(scheduleRequest));
        log.info("Report generation scheduled for cluster: {} with period: {}", scheduleRequest.clusterId(), scheduleRequest.periodMs());
    }

    private void validateClusterId(String clusterId) {
        if (!clusterRepository.existsById(clusterId)) {
            throw new IllegalArgumentException("Cluster with id: " + clusterId + " does not exist.");
        }
    }
}
