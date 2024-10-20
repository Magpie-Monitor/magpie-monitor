package pl.pwr.zpi.metadata.messaging.event.application;

import java.util.List;

public record AggregatedApplicationMetadata(
        Long collectedAtMs,
        String clusterId,
        List<ApplicationMetadata> metadata) {
}
