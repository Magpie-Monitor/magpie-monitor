package pl.pwr.zpi.metadata.dto.cluster;

import lombok.Builder;
import lombok.Data;
import pl.pwr.zpi.reports.enums.Accuracy;

@Data
@Builder
public class ClusterMetadataDTO {
    String clusterId;
    Long updatedAtMillis;
    Accuracy accuracy;
    boolean running;

    public static ClusterMetadataDTO of(String clusterId, boolean running) {
        return ClusterMetadataDTO.builder()
                .clusterId(clusterId)
                .running(running)
                .build();
    }
}
