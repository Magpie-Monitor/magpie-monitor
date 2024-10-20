package pl.pwr.zpi.metadata.messaging.event.node;

import java.util.List;

public record AggregatedNodeMetadata(
        Long collectedAtMs,
        String clusterId,
        List<NodeMetadata> metadata) {
}
