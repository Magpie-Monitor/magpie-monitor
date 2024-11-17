package pl.pwr.zpi.reports.dto.request;

import jakarta.validation.constraints.Min;
import lombok.Builder;
import lombok.NonNull;
import pl.pwr.zpi.cluster.entity.ClusterConfiguration;
import pl.pwr.zpi.reports.dto.report.application.ApplicationConfigurationDTO;
import pl.pwr.zpi.reports.dto.report.node.NodeConfigurationDTO;
import pl.pwr.zpi.reports.enums.Accuracy;

import java.util.List;
@Builder
public record CreateScheduleRequest(
        @NonNull
        String clusterId,
        @NonNull
        @Min(86400000)
        Long periodMs
) {}
