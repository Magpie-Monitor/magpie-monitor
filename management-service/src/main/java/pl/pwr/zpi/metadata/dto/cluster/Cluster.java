package pl.pwr.zpi.metadata.dto.cluster;

import java.util.Objects;

public record Cluster(String clusterId, boolean running) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Cluster cluster = (Cluster) o;
        return Objects.equals(clusterId, cluster.clusterId);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(clusterId);
    }
}
