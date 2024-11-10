package pl.pwr.zpi.metadata.entity;

import pl.pwr.zpi.metadata.dto.application.ApplicationMetadataDTO;
import pl.pwr.zpi.metadata.dto.node.NodeMetadataDTO;
import pl.pwr.zpi.metadata.broker.dto.cluster.ClusterMetadata;

import java.io.Serializable;
import java.util.HashSet;
import java.util.Set;

public record ClusterHistory(
        String id,
        Set<ApplicationMetadataDTO> applicationMetadataDTOS,
        Set<NodeMetadataDTO> nodeMetadataDTOS) implements Serializable {

    public static ClusterHistory of(ClusterMetadata clusterMetadata) {
        return new ClusterHistory(clusterMetadata.clusterId(), new HashSet<>(), new HashSet<>());
    }
}
