package pl.pwr.zpi.metadata.event.dto;

import pl.pwr.zpi.metadata.dto.application.AggregatedApplicationMetadata;
import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;

import java.util.List;

public record ApplicationMetadataUpdated(
        String requestId,
        AggregatedApplicationMetadata metadata) {

    public String clusterId() {
        return metadata.clusterId();
    }

    public List<ApplicationMetadata> applicationMetadata() {
        return metadata.metadata();
    }
}
