package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.messaging.event.node.AggregatedNodeMetadata;

public interface AggregatedNodeMetadataRepository extends MongoRepository<AggregatedNodeMetadata, String> {
}
