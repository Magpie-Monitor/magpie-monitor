package pl.pwr.zpi.reports.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.node.NodeIncident;

import java.util.List;

public interface NodeIncidentRepository extends MongoRepository<NodeIncident, String> {

    List<NodeIncident> findByReportId(String reportId, Pageable pageable);

    Long countByReportId(String reportId);
}
