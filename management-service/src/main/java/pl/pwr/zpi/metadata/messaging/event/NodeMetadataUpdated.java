package pl.pwr.zpi.metadata.messaging.event;

import pl.pwr.zpi.metadata.messaging.event.node.AggregatedNodeMetadata;

import java.util.List;

public record NodeMetadataUpdated(String requestId, List<AggregatedNodeMetadata> metadata) {
}
