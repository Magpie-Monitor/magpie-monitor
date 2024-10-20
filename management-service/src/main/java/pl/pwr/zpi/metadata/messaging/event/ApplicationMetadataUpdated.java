package pl.pwr.zpi.metadata.messaging.event;

import pl.pwr.zpi.metadata.messaging.event.application.AggregatedApplicationMetadata;

import java.util.List;

public record ApplicationMetadataUpdated(String requestId, List<AggregatedApplicationMetadata> metadata) {
}
