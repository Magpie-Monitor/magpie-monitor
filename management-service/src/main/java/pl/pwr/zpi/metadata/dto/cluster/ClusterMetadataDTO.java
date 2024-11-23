package pl.pwr.zpi.metadata.dto.cluster;

import pl.pwr.zpi.reports.enums.Accuracy;

public record ClusterMetadataDTO(
        String clusterId,
        //Long updatedAtMillis,
        //Accuracy accuracy,
        boolean running
) {
}
