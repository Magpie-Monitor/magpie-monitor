package pl.pwr.zpi.metadata.dto.cluster;

import pl.pwr.zpi.metadata.dto.application.ApplicationMetadata;
import pl.pwr.zpi.metadata.dto.node.Node;

import java.io.Serializable;
import java.util.HashSet;
import java.util.Set;

public record ClusterHistory(
        String id,
        Set<ApplicationMetadata> applications,
        Set<Node> nodes) implements Serializable {

    public static ClusterHistory of(ClusterMetadata clusterMetadata) {
        return new ClusterHistory(clusterMetadata.clusterId(), new HashSet<>(), new HashSet<>());
    }
}
