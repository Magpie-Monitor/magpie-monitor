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
import pl.pwr.zpi.reports.enums.Accuracy;
import pl.pwr.zpi.reports.enums.Urgency;
import pl.pwr.zpi.reports.repository.ReportRepository;

import java.util.List;
import java.util.UUID;

@SpringBootApplication
public class MagpieMonitorApplication {

    @Autowired
    public ReportRepository reportRepository;

    @EventListener(ApplicationReadyEvent.class)
    public void init() {
        List<ApplicationReport> appReports = List.of(
                ApplicationReport.builder()
                        .applicationName("test")
                        .accuracy(Accuracy.LOW)
                        .customPrompt("none")
                        .incidents(List.of(
                                ApplicationIncident.builder()
                                        .id(UUID.randomUUID().toString())
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
                    .accuracy(Accuracy.LOW)
                    .customPrompt("none")
                    .incidents(List.of(
                            NodeIncident.builder()
                                    .id(UUID.randomUUID().toString())
                                    .clusterId("test2")
                                    .nodeName("test")
                                    .summary("test")
                                    .build()
                    ))
                    .build()
        );

        System.out.println("hello world, I have just started up");
        Report r = Report.builder()
                .id(UUID.randomUUID().toString())
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
                .scheduledApplicationInsights(null)
                .scheduledNodeInsights(null)
                .build();

        reportRepository.save(r);
    }

    public static void main(String[] args) {
        SpringApplication.run(MagpieMonitorApplication.class, args);
    }
}
