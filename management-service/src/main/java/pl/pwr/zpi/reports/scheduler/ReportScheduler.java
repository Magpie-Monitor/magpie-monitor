package pl.pwr.zpi.reports.scheduler;

import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.TaskScheduler;
import org.springframework.scheduling.concurrent.ThreadPoolTaskScheduler;
import org.springframework.stereotype.Service;
import pl.pwr.zpi.reports.ReportsService;

import java.time.Duration;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

import static com.zaxxer.hikari.util.ClockSource.toMillis;

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

    @PostConstruct
    public void init() {
        reportScheduleRepository.getFirst().ifPresentOrElse(this::generatePastReports, () -> log.info("No report schedule found"));
    }

    private void generatePastReports(ReportSchedule reportSchedule) {
        if (!reportSchedule.getActive()) {
            log.info("Report schedule is inactive, skipping past reports generation.");
            return;
        }
        long breakBetweenReports = Duration.between(reportSchedule.getLastRunTime(), LocalDateTime.now()).toMillis();
        long periodInMillis = reportSchedule.getPeriodInMilis();

        long numberOfReportsToGenerate = breakBetweenReports / periodInMillis;

        LocalDateTime lastRunTime = reportSchedule.getLastRunTime();
        Duration periodDuration = Duration.ofMillis(periodInMillis);

        for (int i = 0; i < numberOfReportsToGenerate; i++) {
            lastRunTime = lastRunTime.plus(periodDuration);
            reportJobService.generateReport(reportSchedule.getLastRunTime(), lastRunTime);
            reportSchedule.setLastRunTime(LocalDateTime.now());
            reportScheduleRepository.save(reportSchedule);
        }

        reportSchedule.setLastRunTime(lastRunTime);
        reportScheduleRepository.save(reportSchedule);

        scheduleNextRun(reportSchedule);
    }

    public void scheduleReport(Long days, Long hours, Long minutes) {
        ReportSchedule reportSchedule = reportScheduleRepository.getFirst()
                .orElseGet(() -> createNewSchedule(days, hours, minutes));

        long periodInMillis = (days * 24 * 60 + hours * 60 + minutes) * 60 * 1000;
        scheduleReport(reportSchedule, periodInMillis);
    }

    public void deactivateSchedule() {
        reportScheduleRepository.getFirst().ifPresent(reportSchedule -> {
            reportSchedule.setActive(false);
            reportScheduleRepository.save(reportSchedule);
            if (scheduledTask != null) {
                scheduledTask.cancel(false);
            }
            log.info("Report schedule deactivated.");
        });
    }

    private void scheduleReport(ReportSchedule reportSchedule, Long periodInMillis) {
        reportSchedule.setPeriodInMilis(periodInMillis);
        reportScheduleRepository.save(reportSchedule);
        scheduleNextRun(reportSchedule);
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
        }, Instant.ofEpochMilli(Math.max(0, delayInMillis)), Duration.ofMillis(reportSchedule.getPeriodInMilis()));
    }

    private ReportSchedule createNewSchedule(Long days, Long hours, Long minutes) {
        return ReportSchedule.builder()
                .lastRunTime(null)
                .periodInMilis((days * 24 * 60 + hours * 60 + minutes) * 60 * 1000)
                .active(true)
                .build();
    }

    public ReportSchedule getLastRunTime() {
        return reportScheduleRepository.getFirst().orElse(null);
    }
}
