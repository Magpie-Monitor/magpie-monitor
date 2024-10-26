package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.dto.node.AggregatedNodeMetadata;

import java.util.Optional;

public interface AggregatedNodeMetadataRepository extends MongoRepository<AggregatedNodeMetadata, String> {
    Optional<AggregatedNodeMetadata> findFirstByClusterIdOrderByCollectedAtMs(String clusterId);
}
