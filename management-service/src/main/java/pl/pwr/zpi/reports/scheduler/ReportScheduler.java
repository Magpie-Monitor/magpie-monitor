package pl.pwr.zpi.reports.scheduler;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.dto.scheduler.ReportSchedule;
import pl.pwr.zpi.reports.repository.ReportScheduleRepository;
import pl.pwr.zpi.reports.service.ReportGenerationService;

import java.util.List;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@AllArgsConstructor
public class ReportScheduler {
    private final ReportScheduleRepository reportScheduleRepository;
    private final ClusterRepository clusterRepository;
    private final ReportGenerationService reportGenerationService;

    @Scheduled(cron = "${report.scheduler.cron}")
    public void generateReports() {
        reportScheduleRepository.findAll().forEach(this::processSchedule);
    }

    private void processSchedule(ReportSchedule schedule) {
        long nextGenerationTime = calculateNextGenerationTime(schedule);

        if (nextGenerationTime > System.currentTimeMillis()) {
            return;
        }

        log.info("Generating report for cluster {}", schedule.getClusterId());

        ClusterConfiguration cluster = clusterRepository.findById(schedule.getClusterId())
                .orElseThrow(() -> new IllegalStateException(
                        String.format("Cluster configuration not found for ID: %s", schedule.getClusterId())
                ));

        CreateReportRequest reportRequest = CreateReportRequest.fromClusterConfiguration(
                cluster,
                schedule.getLastGenerationMs(),
                nextGenerationTime
        );

        reportGenerationService.createReport(reportRequest);

        schedule.setLastGenerationMs(nextGenerationTime);
        reportScheduleRepository.save(schedule);

        log.info("Report generation completed and schedule updated for cluster {}.", schedule.getClusterId());
    }

    private long calculateNextGenerationTime(ReportSchedule schedule) {
        return schedule.getLastGenerationMs() + schedule.getPeriodMs();
    }
}
