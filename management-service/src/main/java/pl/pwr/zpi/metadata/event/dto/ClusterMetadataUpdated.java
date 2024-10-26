package pl.pwr.zpi.metadata.event.dto;

import pl.pwr.zpi.metadata.dto.cluster.AggregatedClusterMetadata;

public record ClusterMetadataUpdated(
        String requestId,
        AggregatedClusterMetadata metadata) {
}
