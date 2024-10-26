package pl.pwr.zpi.metadata.dto.application;

import java.util.List;

public record AggregatedApplicationMetadata(
        Long collectedAtMs,
        String clusterId,
        List<ApplicationMetadata> metadata) {
}
