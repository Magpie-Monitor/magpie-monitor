package pl.pwr.zpi.reports.dto.event;

import lombok.Builder;
import pl.pwr.zpi.reports.dto.request.CreateReportRequest;
import pl.pwr.zpi.reports.enums.Accuracy;

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
            Accuracy accuracy,
            String customPrompt
    ) {
    }

    record NodeConfiguration(
            String nodeName,
            Accuracy accuracy,
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
                                                        configuration.accuracy(),
                                                        configuration.customPrompt()
                                                ))
                                                .toList()
                                )
                                .nodeConfiguration(
                                        reportRequest.nodeConfigurations().stream()
                                                .map(configuration -> new NodeConfiguration(
                                                        configuration.nodeName(),
                                                        configuration.accuracy(),
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
