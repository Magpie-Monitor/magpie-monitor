package pl.pwr.zpi.cluster.dto;

import lombok.NonNull;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;

public record UpdateClusterConfigurationRequest(
        @NonNull
        String id,
        @NonNull
        Accuracy accuracy,
        @NonNull
        Boolean isEnabled,
        @NonNull
        Long generatedEveryMillis,
        @NonNull
        List<Long> slackReceiverIds,
        @NonNull
        List<Long> discordReceiverIds,
        @NonNull
        List<Long> emailReceiverIds,
        @NonNull
        List<ApplicationConfigurationDTO> applicationConfigurations,
        @NonNull
        List<NodeConfigurationDTO> nodeConfigurations
) {

}
