package pl.pwr.zpi.metadata.broker.dto.application;

import java.util.List;

public record AggregatedApplicationMetadata(
        Long collectedAtMs,
        String clusterId,
        List<ApplicationMetadata> metadata) {
}
