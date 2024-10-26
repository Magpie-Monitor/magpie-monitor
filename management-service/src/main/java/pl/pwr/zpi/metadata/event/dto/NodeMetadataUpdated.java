package pl.pwr.zpi.metadata.event.dto;

import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;

public record NodeMetadataUpdated(String requestId, AggregatedNodeMetadata metadata) {
}
