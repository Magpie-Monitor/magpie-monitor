package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.dto.application.AggregatedApplicationMetadata;

import java.util.Optional;

public interface AggregatedApplicationMetadataRepository extends MongoRepository<AggregatedApplicationMetadata, String> {
    Optional<AggregatedApplicationMetadata> findFirstByClusterIdOrderByCollectedAtMs(String clusterId);
}