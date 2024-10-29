package pl.pwr.zpi.metadata.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.metadata.dto.cluster.ClusterHistory;

public interface ClusterHistoryRepository extends MongoRepository<ClusterHistory, String> {
}
