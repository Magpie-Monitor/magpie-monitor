package pl.pwr.zpi.reports.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncidentSource;

import java.util.List;

public interface ApplicationIncidentSourcesRepository extends MongoRepository<ApplicationIncidentSource, String> {
    List<ApplicationIncidentSource> findByIncidentId(String incidentId, Pageable pageable);
    Long countByIncidentId(String incidentId);
}
