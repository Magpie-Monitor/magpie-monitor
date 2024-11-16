package pl.pwr.zpi.reports.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.node.NodeIncidentSource;

import java.util.List;

public interface NodeIncidentSourcesRepository extends MongoRepository<NodeIncidentSource, String> {
    List<NodeIncidentSource> findByIncidentId(String incidentId, Pageable pageable);
    Long countByIncidentId(String incidentId);
}
