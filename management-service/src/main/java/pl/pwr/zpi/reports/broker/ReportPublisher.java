package pl.pwr.zpi.reports.broker;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.reports.dto.event.ReportRequested;

@Slf4j
@Component
@RequiredArgsConstructor
public class ReportPublisher {

    @Value("${kafka.report.requested.topic}")
    private String REPORT_REQUESTED_TOPIC;
    private final KafkaTemplate<String, ReportRequested> kafkaTemplate;

    public void publishReportRequestedEvent(ReportRequested reportRequested) {
        kafkaTemplate.send(
                REPORT_REQUESTED_TOPIC,
                reportRequested.correlationId(),
                reportRequested
        ).whenComplete((result, ex) -> {
            if (ex != null) {
                log.error("error publishing report requested event: {}", ex.getMessage());
            }
        });
    }
}
