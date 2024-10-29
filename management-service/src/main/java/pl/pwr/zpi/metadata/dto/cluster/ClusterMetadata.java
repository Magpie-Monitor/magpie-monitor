package pl.pwr.zpi.metadata.dto.cluster;

import java.io.Serializable;

public record ClusterMetadata(
        String name,
        boolean running) implements Serializable {
}
