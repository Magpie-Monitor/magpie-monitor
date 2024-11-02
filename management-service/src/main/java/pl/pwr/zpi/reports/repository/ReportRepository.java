package pl.pwr.zpi.reports.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.repository.projection.*;

import java.util.List;
import java.util.Optional;

public interface ReportRepository extends MongoRepository<Report, String> {

    List<ReportSummaryProjection> findAllProjectedBy();

    Optional<ReportDetailedSummaryProjection> findProjectedBy(String reportId);

    Optional<ReportIncidentsProjection> findProjectedById(String reportId);
}
