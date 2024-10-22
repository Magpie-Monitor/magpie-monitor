package pl.pwr.zpi.metadata.messaging.event;

import pl.pwr.zpi.metadata.messaging.event.node.AggregatedNodeMetadata;

public record NodeMetadataUpdated(String requestId, AggregatedNodeMetadata metadata) {
}
