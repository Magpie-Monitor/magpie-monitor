package pl.pwr.zpi.reports.dto.request;

import pl.pwr.zpi.reports.dto.report.application.ApplicationConfiguration;
import pl.pwr.zpi.reports.dto.report.node.NodeConfiguration;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

public record CreateReportRequest(
        String clusterId,
        Accuracy accuracy, // TODO - validate if needed with reports service
        Long sinceMs,
        Long toMs,
        List<Long> slackReceiverIds,
        List<Long> discordReceiverIds,
        List<Long> mailReceiverIds,
        List<ApplicationConfiguration> applicationConfigurations,
        List<NodeConfiguration> nodeConfigurations
) {
}
