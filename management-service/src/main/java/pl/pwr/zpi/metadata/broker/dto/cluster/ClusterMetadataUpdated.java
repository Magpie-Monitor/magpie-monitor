package pl.pwr.zpi.metadata.broker.dto.cluster;

import java.util.List;

public record ClusterMetadataUpdated(
        String correlationId,
        AggregatedClusterMetadata metadata) {

    public List<ClusterMetadata> clusterMetadata() {
        return metadata.metadata();
    }
}
