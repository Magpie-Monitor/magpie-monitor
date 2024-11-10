package pl.pwr.zpi.reports.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;

import java.util.List;

public interface ApplicationIncidentRepository extends MongoRepository<ApplicationIncident, String> {
    List<ApplicationIncident> findByReportId(String reportId, Pageable pageable);

    Long countByReportId(String reportId);
}
