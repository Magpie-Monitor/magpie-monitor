package pl.pwr.zpi.metadata.event.dto;

import pl.pwr.zpi.metadata.dto.cluster.AggregatedClusterMetadata;
import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;

import java.util.List;

public record ClusterMetadataUpdated(
        String requestId,
        AggregatedClusterMetadata metadata) {

    public List<ClusterMetadata> clusterMetadata() {
        return metadata.metadata();
    }
}
