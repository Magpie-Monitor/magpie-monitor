package pl.pwr.zpi.metadata.dto.cluster;

import java.util.List;

public record AggregatedClusterMetadata(
        Long collectedAtMs,
        List<ClusterMetadata> metadata
) {
}
