package pl.pwr.zpi.cluster;

import lombok.Builder;

@Builder
public record ClusterConfiguration(
        String id,
        boolean running,
        String precision
) {
}
