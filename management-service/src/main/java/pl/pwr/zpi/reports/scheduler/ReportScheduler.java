package pl.pwr.zpi.reports.scheduler;

import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.TaskScheduler;
import org.springframework.scheduling.concurrent.ThreadPoolTaskScheduler;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.ReportsService;

import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.concurrent.ScheduledFuture;

@Service
@Slf4j
public class ReportScheduler {

    private final TaskScheduler taskScheduler;
    private final ReportsService reportJobService;
    private final ReportScheduleRepository reportScheduleRepository;
    private ScheduledFuture<?> scheduledTask;

    @Autowired
    public ReportScheduler(ReportsService reportJobService, ReportScheduleRepository reportScheduleRepository) {
        this.reportJobService = reportJobService;
        this.reportScheduleRepository = reportScheduleRepository;
        this.taskScheduler = new ThreadPoolTaskScheduler();
        ((ThreadPoolTaskScheduler) taskScheduler).initialize();
    }

//    @PostConstruct
    public void init() {
        reportScheduleRepository.getFirst().ifPresentOrElse(this::scheduleNextRun, () -> log.info("No report schedule found"));
    }

    public void scheduleReport(Long days, Long hours, Long minutes) {
        ReportSchedule reportSchedule = reportScheduleRepository.getFirst()
                .orElseGet(() -> createNewSchedule(days, hours, minutes));

        long periodInMillis = (days * 24 * 60 + hours * 60 + minutes) * 60 * 1000;
        reportSchedule.setPeriodInMilis(periodInMillis);

        reportScheduleRepository.save(reportSchedule);

        handleMissedRuns(reportSchedule);

        scheduleNextRun(reportSchedule);
    }

    private void handleMissedRuns(ReportSchedule reportSchedule) {
        if (reportSchedule.getLastRunTime() != null) {
            LocalDateTime now = LocalDateTime.now();
            long periodInMillis = reportSchedule.getPeriodInMilis();
            LocalDateTime lastRunTime = reportSchedule.getLastRunTime();
            Duration periodDuration = Duration.ofMillis(periodInMillis);

            while (lastRunTime.plus(periodDuration).isBefore(now)) {
                lastRunTime = lastRunTime.plus(periodDuration);
                reportJobService.generateReport(reportSchedule.getLastRunTime(), lastRunTime);
            }

            reportSchedule.setLastRunTime(lastRunTime);
            reportScheduleRepository.save(reportSchedule);
        }
    }

    private void scheduleNextRun(ReportSchedule reportSchedule) {
        if (scheduledTask != null) {
            scheduledTask.cancel(false);
        }

        long lastRunMillis = (reportSchedule.getLastRunTime() == null) ? 0
                : System.currentTimeMillis() - reportSchedule.getLastRunTime().atZone(ZoneId.systemDefault()).toInstant().toEpochMilli();
        long delayInMillis = reportSchedule.getPeriodInMilis() - lastRunMillis;

        scheduledTask = taskScheduler.scheduleAtFixedRate(() -> {
            reportJobService.generateReport(reportSchedule.getLastRunTime(), LocalDateTime.now());

            reportSchedule.setLastRunTime(LocalDateTime.now());
            reportScheduleRepository.save(reportSchedule);
        }, Duration.ofMillis(Math.max(0, delayInMillis)));
    }

    private ReportSchedule createNewSchedule(Long days, Long hours, Long minutes) {
        ReportSchedule reportSchedule = new ReportSchedule();
        reportSchedule.setLastRunTime(null);
        reportSchedule.setPeriodInMilis((days * 24 * 60 + hours * 60 + minutes) * 60 * 1000);
        return reportSchedule;
    }
}
