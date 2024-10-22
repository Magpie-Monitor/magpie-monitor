package pl.pwr.zpi.metadata.messaging.event;

import pl.pwr.zpi.metadata.messaging.event.application.AggregatedApplicationMetadata;

public record ApplicationMetadataUpdated(String requestId, AggregatedApplicationMetadata metadata) {
}
