package pl.pwr.zpi.metadata.event.dto.node;

import java.util.List;

public record AggregatedNodeMetadata(
        Long correlationId,
        String clusterId,
        List<NodeMetadata> metadata) {
}
