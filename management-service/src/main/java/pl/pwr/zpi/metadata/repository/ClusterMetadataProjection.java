package pl.pwr.zpi.metadata.repository;

import pl.pwr.zpi.metadata.dto.cluster.ClusterMetadata;

import java.util.List;

public interface ClusterMetadataProjection {

    List<ClusterMetadata> getMetadata();
}
