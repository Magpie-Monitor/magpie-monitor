package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.messaging.event.application.AggregatedApplicationMetadata;

public interface AggregatedApplicationMetadataRepository extends MongoRepository<AggregatedApplicationMetadata, String> {
}
