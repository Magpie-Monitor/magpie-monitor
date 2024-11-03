package pl.pwr.zpi.metadata.broker.dto.node;

import java.util.List;

public record NodeMetadataUpdated(String correlationId, AggregatedNodeMetadata metadata) {

    public String clusterId() {
        return metadata.clusterId();
    }

    public List<NodeMetadata> nodeMetadata() {
        return metadata.metadata();
    }
}
