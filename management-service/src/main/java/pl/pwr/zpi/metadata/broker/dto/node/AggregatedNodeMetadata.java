package pl.pwr.zpi.metadata.broker.dto.node;

import java.util.List;

public record AggregatedNodeMetadata(
        Long collectedAtMs,
        String clusterId,
        List<NodeMetadata> metadata) {
}
