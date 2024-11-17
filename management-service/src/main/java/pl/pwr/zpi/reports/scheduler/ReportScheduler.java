package pl.pwr.zpi.reports.scheduler;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.cluster.repository.ClusterRepository;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.dto.scheduler.ClusterSchedule;
import pl.pwr.zpi.reports.repository.SchedulerRepository;
import pl.pwr.zpi.reports.service.ReportGenerationService;

import java.util.List;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@AllArgsConstructor
public class ReportScheduler {
    private final SchedulerRepository schedulerRepository;
    private final ClusterRepository clusterRepository;
    private final ReportGenerationService reportGenerationService;

    private static final long TEN_MINUTES_IN_MILLIS = TimeUnit.MINUTES.toMillis(10);

    @Scheduled(cron = "0 */10 * * * *")
    public void generateReports() {
        long currentTime = System.currentTimeMillis();

        List<ClusterSchedule> schedules = schedulerRepository.findAll();

        schedules.forEach(schedule -> {
            try {
                processSchedule(schedule, currentTime);
            } catch (Exception e) {
                log.error("Error processing schedule for cluster {}: {}", schedule.getClusterId(), e.getMessage(), e);
            }
        });
    }

    private void processSchedule(ClusterSchedule schedule, long currentTime) {
        long nextGenerationTime = schedule.getLastGenerationMs() + schedule.getPeriodMs();

        if (nextGenerationTime > currentTime + TEN_MINUTES_IN_MILLIS) {
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
        schedulerRepository.save(schedule);

        log.info("Report generation completed and schedule updated for cluster {}.", schedule.getClusterId());
    }
}
