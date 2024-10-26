package pl.pwr.zpi.metadata.dto.node;

import java.util.List;

public record AggregatedNodeMetadata(
        Long collectedAtMs,
        String clusterId,
        List<NodeMetadata> metadata) {
}
