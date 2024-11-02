package pl.pwr.zpi.reports.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.application.ApplicationIncident;

public interface ApplicationIncidentRepository extends MongoRepository<ApplicationIncident, String> {
}
