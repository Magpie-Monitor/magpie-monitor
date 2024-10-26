package pl.pwr.zpi.metadata.event.dto;

import pl.pwr.zpi.metadata.dto.application.AggregatedApplicationMetadata;

public record ApplicationMetadataUpdated(
        String requestId,
        AggregatedApplicationMetadata metadata) {
}
