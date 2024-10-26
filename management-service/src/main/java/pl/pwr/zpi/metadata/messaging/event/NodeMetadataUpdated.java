package pl.pwr.zpi.metadata.messaging.event;

import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;

public record NodeMetadataUpdated(String requestId, AggregatedNodeMetadata metadata) {
}
