package pl.pwr.zpi.reports.broker;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import pl.pwr.zpi.reports.dto.event.ReportGenerated;
import pl.pwr.zpi.reports.dto.event.ReportRequestFailed;
import pl.pwr.zpi.reports.service.ReportGenerationService;
import pl.pwr.zpi.utils.mapper.JsonMapper;

@Slf4j
@Component
@RequiredArgsConstructor
public class ReportConsumer {

    private final JsonMapper mapper;
    private final ReportGenerationService reportGenerationService;

    @KafkaListener(topics = "${kafka.report.generated.topic}")
    public void listenForReportGeneratedEvent(String message) {
        ReportGenerated report = mapper.fromJson(message, ReportGenerated.class);
        log.info("Received report created event: {}", report);
        reportGenerationService.handleReportGenerated(report);
    }

    @KafkaListener(topics = "${kafka.report.request.failed.topic}")
    public void listenForReportRequestFailedEvent(String message) {
        ReportRequestFailed request = mapper.fromJson(message, ReportRequestFailed.class);
        log.info("Received report request failed: {}", request);
        reportGenerationService.handleReportGenerationError(request);
    }
}
