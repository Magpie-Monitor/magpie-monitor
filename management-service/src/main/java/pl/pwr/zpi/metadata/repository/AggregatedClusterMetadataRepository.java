package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.dto.cluster.AggregatedClusterMetadata;

import java.util.Optional;

public interface AggregatedClusterMetadataRepository extends MongoRepository<AggregatedClusterMetadata, String> {
    Optional<AggregatedClusterMetadata> findFirstByMetadataClusterId(String clusterId);

    Optional<AggregatedClusterMetadata> findFirstByOrderByCollectedAtMsDesc();

    boolean existsByMetadataClusterId(String clusterId);
}
