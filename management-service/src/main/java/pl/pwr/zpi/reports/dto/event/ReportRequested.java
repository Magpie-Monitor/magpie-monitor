package pl.pwr.zpi.reports.dto.event;

import lombok.Builder;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.enums.Precision;

import java.util.List;
import java.util.UUID;

@Builder
public record ReportRequested(
        String correlationId,
        ReportRequest reportRequest
) {

    @Builder
    record ReportRequest(
            String clusterId,
            Long sinceMs,
            Long toMs,
            List<ApplicationConfiguration> applicationConfiguration,
            List<NodeConfiguration> nodeConfiguration,
            Integer maxLength
    ) {
    }

    record ApplicationConfiguration(
            String applicationName,
            Precision precision,
            String customPrompt
    ) {
    }

    record NodeConfiguration(
            String nodeName,
            Precision precision,
            String customPrompt
    ) {
    }

    public static ReportRequested of(CreateReportRequest reportRequest) {
        return ReportRequested.builder()
                .correlationId(UUID.randomUUID().toString())
                .reportRequest(
                        ReportRequest.builder()
                                .clusterId(reportRequest.clusterId())
                                .sinceMs(reportRequest.sinceMs())
                                .toMs(reportRequest.toMs())
                                .applicationConfiguration(
                                        reportRequest.applicationConfigurations().stream()
                                                .map(configuration -> new ApplicationConfiguration(
                                                        configuration.applicationName(),
                                                        configuration.precision(),
                                                        configuration.customPrompt()
                                                ))
                                                .toList()
                                )
                                .nodeConfiguration(
                                        reportRequest.nodeConfigurations().stream()
                                                .map(configuration -> new NodeConfiguration(
                                                        configuration.nodeName(),
                                                        configuration.precision(),
                                                        configuration.customPrompt()
                                                ))
                                                .toList()
                                )
                                .maxLength(100)
                                .build()
                )
                .build();
    }
}
