package pl.pwr.zpi;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;
import pl.pwr.zpi.reports.entity.report.application.ApplicationReport;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;
import pl.pwr.zpi.reports.entity.report.node.NodeReport;
import pl.pwr.zpi.reports.entity.report.node.scheduled.ScheduledNodeInsight;
import pl.pwr.zpi.reports.enums.Precision;
import pl.pwr.zpi.reports.enums.Urgency;
import pl.pwr.zpi.reports.repository.ReportRepository;

import java.util.List;

@SpringBootApplication
public class MagpieMonitorApplication {

    @Autowired
    public ReportRepository reportRepository;

    @EventListener(ApplicationReadyEvent.class)
    public void doSomethingAfterStartup() {
        List<ApplicationReport> appReports = List.of(
                ApplicationReport.builder()
                        .applicationName("test")
                        .precision(Precision.LOW)
                        .customPrompt("none")
                        .incidents(List.of(
                                ApplicationIncident.builder()
                                        .id("test")
                                        .clusterId("test2")
                                        .applicationName("test")
                                        .category("test")
                                        .summary("test")
                                        .recommendation("test")
                                        .build()
                        ))
                        .build()
        );

        List<NodeReport> nodeReports = List.of(
            NodeReport.builder()
                    .node("test")
                    .precision(Precision.LOW)
                    .customPrompt("none")
                    .nodeIncidents(List.of(
                            NodeIncident.builder()
                                    .id("test2")
                                    .clusterId("test2")
                                    .nodeName("test")
                                    .summary("test")
                                    .build()
                    ))
                    .build()
        );

        System.out.println("hello world, I have just started up");
        Report r = Report.builder()
                .clusterId("test")
                .sinceMs(3000L)
                .toMs(5000L)
                .requestedAtMs(5000L)
                .title("test")
                .nodeReports(nodeReports)
                .applicationReports(appReports)
                .totalApplicationEntries(100)
                .totalNodeEntries(100)
                .urgency(Urgency.MEDIUM)
                .scheduledApplicationInsights(List.of())
                .scheduledNodeInsights(List.of())
                .build();

        reportRepository.save(r);
    }

    public static void main(String[] args) {
        SpringApplication.run(MagpieMonitorApplication.class, args);
    }
}
