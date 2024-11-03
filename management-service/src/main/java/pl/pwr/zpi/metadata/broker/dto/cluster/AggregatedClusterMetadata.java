package pl.pwr.zpi.metadata.broker.dto.cluster;

import java.util.List;

public record AggregatedClusterMetadata(
        Long collectedAtMs,
        List<ClusterMetadata> metadata
) {
}
