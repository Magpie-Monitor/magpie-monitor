package pl.pwr.zpi.metadata.entity;

import pl.pwr.zpi.metadata.dto.application.Application;
import pl.pwr.zpi.metadata.dto.node.Node;
import pl.pwr.zpi.metadata.event.dto.cluster.ClusterMetadata;

import java.io.Serializable;
import java.util.HashSet;
import java.util.Set;

public record ClusterHistory(
        String id,
        Set<Application> applications,
        Set<Node> nodes) implements Serializable {

    public static ClusterHistory of(ClusterMetadata clusterMetadata) {
        return new ClusterHistory(clusterMetadata.clusterId(), new HashSet<>(), new HashSet<>());
    }
}
