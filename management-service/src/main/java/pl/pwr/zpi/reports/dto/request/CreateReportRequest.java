package pl.pwr.zpi.reports.dto.request;

import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

public record CreateReportRequest(
        String clusterId,
        Accuracy accuracy,
        Long sinceMs,
        Long toMs,
        List<Long> slackReceiverIds,
        List<Long> discordReceiverIds,
        List<Long> mailReceiverIds,
        List<ApplicationConfiguration> applicationConfigurations,
        List<NodeConfiguration> nodeConfigurations
) {

    public record ApplicationConfiguration(
            String applicationName,
            String customPrompt,
            boolean enabled,
            Accuracy accuracy
    ) {
    }

    public record NodeConfiguration(
            String nodeName,
            String customPrompt,
            boolean enabled,
            Accuracy accuracy
    ) {
    }
}