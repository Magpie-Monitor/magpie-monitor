package pl.pwr.zpi.metadata.broker.dto.application;

import java.util.List;

public record ApplicationMetadataUpdated(
        String correlationId,
        AggregatedApplicationMetadata metadata) {

    public String clusterId() {
        return metadata.clusterId();
    }

    public List<ApplicationMetadata> applicationMetadata() {
        return metadata.metadata();
    }
}
