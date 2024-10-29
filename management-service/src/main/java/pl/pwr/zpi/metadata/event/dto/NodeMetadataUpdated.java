package pl.pwr.zpi.metadata.event.dto;

import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;
import pl.pwr.zpi.metadata.dto.node.NodeMetadata;

import java.util.List;

public record NodeMetadataUpdated(String requestId, AggregatedNodeMetadata metadata) {

    public String clusterId() {
        return metadata.clusterId();
    }

    public List<NodeMetadata> nodeMetadata() {
        return metadata.metadata();
    }
}
