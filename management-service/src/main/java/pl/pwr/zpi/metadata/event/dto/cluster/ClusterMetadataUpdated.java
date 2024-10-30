package pl.pwr.zpi.metadata.event.dto.cluster;

import java.util.List;

public record ClusterMetadataUpdated(
        String correlationId,
        AggregatedClusterMetadata metadata) {

    public List<ClusterMetadata> clusterMetadata() {
        return metadata.metadata();
    }
}
