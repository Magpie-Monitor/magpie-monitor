package pl.pwr.zpi.reports.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;

public interface NodeIncidentRepository extends MongoRepository<NodeIncident, String> {
}
