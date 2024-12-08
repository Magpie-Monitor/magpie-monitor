package pl.pwr.zpi.reports.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import pl.pwr.zpi.reports.entity.report.Report;
import pl.pwr.zpi.reports.enums.ReportType;
import pl.pwr.zpi.reports.repository.projection.ReportDetailedSummaryProjection;
import pl.pwr.zpi.reports.repository.projection.ReportIncidentsProjection;
import pl.pwr.zpi.reports.repository.projection.ReportRequestedAtMillisProjection;
import pl.pwr.zpi.reports.repository.projection.ReportSummaryProjection;

import java.util.List;
import java.util.Optional;

public interface ReportRepository extends MongoRepository<Report, String> {

    List<ReportSummaryProjection> findAllByReportType(ReportType reportType);

    Optional<ReportDetailedSummaryProjection> findProjectedDetailedById(String reportId);

    Optional<ReportIncidentsProjection> findProjectedIncidentsById(String reportId);
    Optional<ReportDetailedSummaryProjection> findFirstByOrderByRequestedAtMsDesc();
}
