package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.messaging.event.node.AggregatedNodeMetadata;

import java.util.List;

public interface AggregatedNodeMetadataRepository extends MongoRepository<AggregatedNodeMetadata, String> {
    List<AggregatedNodeMetadata> findFirstByClusterIdAndOrderByCollectedAtMs(String clusterId);
}
